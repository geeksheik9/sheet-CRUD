package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	log "github.com/sirupsen/logrus"

	model "github.com/geeksheik9/sheet-CRUD/models"
)

//key is intentionally unexported so that there are no collisions with other packages
type key int

const (
	//HTTPHeaderAPIKey is the header for the rbac apikey.
	HTTPHeaderAPIKey string = "apikey"
	//HTTPHeaderAuthorization is the header JWT bearer token
	HTTPHeaderAuthorization = "Authorization"
	//DefaultTimeout is the default timeout for api clients
	DefaultTimeout = 15 * time.Second
)

const (
	// Note: context keys must have unique values.
	//contextAPIKey has an arbitrary value of zero. Any value is okay.
	contextAPIKey key = iota
	//contextIncludeObsolete has an arbitrary value of one. Any value is okay.
	contextIncludeObsolete
)

//APIClient base apiClient interface
type APIClient interface {
	Get(ctx context.Context, path string, queryParams *url.Values) (*model.APIResponse, error)
	Do(ctx context.Context, request *http.Request) (*model.APIResponse, error)
	Post(ctx context.Context, path string, body io.Reader) (*model.APIResponse, error)
	Put(ctx context.Context, path string, body io.Reader) (*model.APIResponse, error)
	Delete(ctx context.Context, path string, body io.Reader) (*model.APIResponse, error)
}

//Client holds the baseURL, client and userAgent.
type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	HTTPClient *http.Client
	// some APIs don't require JWT token/authentication, this flag indicates to the client when its API doesn't require
	// credentials. Default is "true"
	RequiresAuthentication bool
}

type httpResult struct {
	Raw   json.RawMessage
	Code  int
	Error error
}

// GetHTTPClient returns http client initialized with the input certificate pool
func GetHTTPClient(certPool *x509.CertPool) *http.Client {
	if certPool == nil {
		httpClient := http.DefaultClient
		httpClient.Timeout = DefaultTimeout
		return httpClient
	}
	return &http.Client{Timeout: DefaultTimeout, Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: certPool}}}
}

//InitClient inits the client given the params passed in.
func InitClient(baseURL string, userAgent string, httpClient *http.Client, timeout time.Duration) (*Client, error) {
	baseEndpoint, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if timeout > 0 {
		httpClient.Timeout = timeout * time.Second
	} else {
		//default to 15 seconds
		httpClient.Timeout = DefaultTimeout
	}
	c := &Client{BaseURL: baseEndpoint,
		UserAgent:              userAgent,
		HTTPClient:             httpClient,
		RequiresAuthentication: true,
	}

	return c, nil
}

//Get basic HTTP get call with support for request parameters and query parameters
func (c *Client) Get(ctx context.Context, urlPath string, queryParams *url.Values) (*model.APIResponse, error) {
	var u url.URL
	if queryParams != nil && len(*queryParams) > 0 {
		q := c.BaseURL.Query()
		for key, value := range *queryParams {
			q[key] = append([]string{}, value...)
		}
		u = *c.BaseURL
		u.Path = path.Join(u.Path, urlPath)
		u.RawQuery = q.Encode()
	} else {
		u = *c.BaseURL
		u.Path = path.Join(u.Path, urlPath)
	}

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Debugf("Failed to : %s\n", err)
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", c.UserAgent)
	return c.Do(ctx, request)
}

//Do executes a HTTP request
func (c *Client) Do(ctx context.Context, request *http.Request) (*model.APIResponse, error) {
	if request != nil { // some unit tests use a nil request
		log.Debugf("apiclient Do(): method %v, url %v", request.Method, request.URL)
	}

	var resp = &model.APIResponse{}

	if c.RequiresAuthentication {
		apiKey, err := APIKeyFromContext(ctx)
		if err != nil {
			resp.StatusCode = http.StatusBadRequest
			return resp, err
		}
		request.Header.Set(HTTPHeaderAuthorization, "Bearer "+apiKey)
	}

	request = request.WithContext(ctx)

	var response *http.Response
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		log.Debugf("Error sending HTTP request to %s: %v\n", request.URL, err)
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return resp, ctx.Err()
		default:
		}

		return resp, err
	}
	defer response.Body.Close()
	if response.Body != nil {
		resp.Body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			resp.StatusCode = http.StatusInternalServerError
			return resp, err
		}
	}
	resp.StatusCode = response.StatusCode
	return resp, err
}

//Put creates a put request and calls Do
func (c *Client) Put(ctx context.Context, urlPath string, body io.Reader) (*model.APIResponse, error) {
	u := *c.BaseURL
	u.Path = path.Join(u.Path, urlPath)
	request, err := http.NewRequest(http.MethodPut, u.String(), body)
	if err != nil {
		log.Debugf("Failed to create PUT : %v\n", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", c.UserAgent)
	return c.Do(ctx, request)
}

//Delete creates a Delete request and calls Do
func (c *Client) Delete(ctx context.Context, urlPath string, body io.Reader) (*model.APIResponse, error) {
	u := *c.BaseURL
	u.Path = path.Join(u.Path, urlPath)
	request, err := http.NewRequest(http.MethodDelete, u.String(), body)
	if err != nil {
		log.Debugf("Failed to create DELETE : %v\n", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", c.UserAgent)
	return c.Do(ctx, request)
}

//Post creates a post request and calls Do
func (c *Client) Post(ctx context.Context, urlPath string, body io.Reader) (*model.APIResponse, error) {
	u := *c.BaseURL
	u.Path = path.Join(u.Path, urlPath)
	request, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		log.Debugf("Failed to create POST : %v\n", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", c.UserAgent)
	return c.Do(ctx, request)
}

//IsSuccessful checks if the code is a 200 code.
func IsSuccessful(s int) bool {
	return s/100 == 2
}

// WithAPIKey puts the JWT token into the current context.
func WithAPIKey(ctx context.Context, jwtToken string) context.Context {
	return context.WithValue(ctx, contextAPIKey, jwtToken)
}

// APIKeyFromContext returns the ApiKey from the context.
func APIKeyFromContext(ctx context.Context) (string, error) {
	v := ctx.Value(contextAPIKey)
	if v == nil {
		return "", ErrMissingAPIKeyInContext
	}
	return v.(string), nil
}

// ErrMissingAPIKeyInContext is the error returned when the context does not contain ApiKey
var ErrMissingAPIKeyInContext = errors.New("missing ApiKey in context")

// IncludeObsoleteQueryParam is the query parameter key for include obsolete
const IncludeObsoleteQueryParam = "includeDeleted"

// PageCountQueryParam is the query parameter key for page count
const PageCountQueryParam = "pageCount"

// PageNumberQueryParam is the query parameter key for page number
const PageNumberQueryParam = "pageNumber"

func Handler(resp *model.APIResponse, err error) *httpResult {
	result := httpResult{}
	result.Error = err
	if resp != nil {
		result.Code = resp.StatusCode
	} else {
		log.Error("apiclient handler: API response was null")
		result.Code = http.StatusInternalServerError
		return &result
	}

	if err != nil {
		log.Errorf("apiclient handler: API request error:  %+v", err)
		return &result
	}

	if !IsSuccessful(resp.StatusCode) {
		log.Errorf("apiclient handler: API request was unsuccessful:  status code %d", result.Code)
		if resp.Body != nil {
			errResp := ErrorResponse{}
			err = json.Unmarshal(resp.Body, &errResp)
			if err != nil {
				// return unmarshal error
				log.Errorf("apiclient handler unmarshaling error: %v", err)
				result.Error = err
				return &result
			}
			// else pass original API error to the caller
			result.Error = errors.New(errResp.Error)
		}
		return &result
	}

	// happy path
	result.Raw = resp.Body
	return &result
}
