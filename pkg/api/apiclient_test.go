package api

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	model "github.com/geeksheik9/sheet-CRUD/models"
)

const (
	baseURLSuffix = "v2"
)

func TestApitClient_GetHttpClient_CertNotNil(t *testing.T) {
	certPool := x509.NewCertPool()
	client := GetHTTPClient(certPool)
	if client == nil {
		t.Error("Failed to intialize client")
	}
}

func TestApitClient_GetHttpClient_CertNil(t *testing.T) {
	client := GetHTTPClient(nil)
	if client == nil {
		t.Error("Failed to intialize client")
	}
}

func TestApiClient_InitClient(t *testing.T) {
	c, _ := InitClient("http://test", "test", nil, 0)
	if c.UserAgent != "test" {
		t.Errorf("Init client wrong user agent, expecting %v, got %v", "test", c.UserAgent)
	}
}

func TestApiClient_Do_noKey(t *testing.T) {
	c, _ := InitClient("http://test", "test", nil, 15)
	_, err := c.Do(context.Background(), nil)

	if err == nil {
		t.Errorf("Init client invalid key error, expecting %v, got %v", "test", c.UserAgent)
	}
}

func TestApiClient_Do_httpInternalServerError(t *testing.T) {
	c, _ := InitClient("http://test", "test", http.DefaultClient, 15)
	c.RequiresAuthentication = false
	request, err := http.NewRequest(http.MethodGet, "http://test", nil)
	if err != nil {
		t.Errorf("Do unable to create request")
	}
	request.URL = nil
	resp, err := c.Do(context.Background(), request)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Do expected http interal server error")
	}
}

func TestApiClient_Get_VerifyRequestURL(t *testing.T) {

	executeRequest := func(ctx context.Context, client APIClient, path string, queryParams *url.Values, body io.Reader) (*model.APIResponse, error) {
		resp, err := client.Get(ctx, path, queryParams)
		// intentionally return response and error objects As-Is
		return resp, err
	}

	v := url.Values{}
	v.Set("name", "Ava")
	v.Set("friend", "Jess")

	verifyRequestURL(t, "", nil, executeRequest)
	verifyRequestURL(t, "", &v, executeRequest)
	verifyRequestURL(t, baseURLSuffix, nil, executeRequest)
	verifyRequestURL(t, baseURLSuffix, &v, executeRequest)
}

func TestApiClient_Post_VerifyRequestURL(t *testing.T) {

	executeRequest := func(ctx context.Context, client APIClient, path string, queryParams *url.Values, body io.Reader) (*model.APIResponse, error) {
		resp, err := client.Post(ctx, path, body)
		// intentionally return response and error objects As-Is
		return resp, err
	}

	verifyRequestURL(t, "", nil, executeRequest)
	verifyRequestURL(t, baseURLSuffix, nil, executeRequest)
}

func TestApiClient_Put_VerifyRequestURL(t *testing.T) {

	executeRequest := func(ctx context.Context, client APIClient, path string, queryParams *url.Values, body io.Reader) (*model.APIResponse, error) {
		resp, err := client.Put(ctx, path, body)
		// intentionally return response and error objects As-Is
		return resp, err
	}

	verifyRequestURL(t, "", nil, executeRequest)
	verifyRequestURL(t, baseURLSuffix, nil, executeRequest)
}

func TestApiClient_Delete_VerifyRequestURL(t *testing.T) {

	executeRequest := func(ctx context.Context, client APIClient, path string, queryParams *url.Values, body io.Reader) (*model.APIResponse, error) {
		resp, err := client.Delete(ctx, path, body)
		// intentionally return response and error objects As-Is
		return resp, err
	}

	verifyRequestURL(t, "", nil, executeRequest)
	verifyRequestURL(t, baseURLSuffix, nil, executeRequest)
}

func TestAPIClient_handler_nilResponse(t *testing.T) {
	result := Handler(nil, nil)
	if result.Code != http.StatusInternalServerError {
		t.Errorf("Error testing handler\n\tExpected:\t%v\n\tReceived:\t%v", http.StatusInternalServerError, result)
	}
}

func TestAPIClient_handler_ExpectedError(t *testing.T) {
	model := model.APIResponse{}
	result := Handler(&model, errors.New("this is an error"))
	if result.Error.Error() != "this is an error" {
		t.Errorf("Error testing handler\n\tExpected\tthis is an error\n\tReceived:\t%v", result.Error)
	}
}

func TestAPIClient_handler_NotSuccessful_EmptyBody(t *testing.T) {
	model := model.APIResponse{
		StatusCode: 403,
	}
	result := Handler(&model, nil)
	if result.Code != 403 {
		t.Errorf("Error testing handler\n\tExpected:\t400\n\tReceived:\t%v", result.Code)
	}
}

func TestAPIClient_handler_NotSuccessful_Error(t *testing.T) {
	model := model.APIResponse{
		Body:       []byte("this is a byte slice"),
		StatusCode: 400,
	}
	result := Handler(&model, nil)
	if result.Error == nil {
		t.Errorf("Error testing handler\n\tExpected\t<not nil>\n\tReceived:\t%v", result.Error)
	}
}

type makeRequest func(ctx context.Context, client APIClient, path string, queryParams *url.Values, body io.Reader) (*model.APIResponse, error)

func verifyRequestURL(t *testing.T, baseURLSuffix string, queryParams *url.Values, request makeRequest) {
	URLSuffix := "roles/bus"
	var uri string
	if baseURLSuffix == "" {
		uri = "/" + URLSuffix
	} else {
		uri = path.Join("/", baseURLSuffix, URLSuffix)
	}
	if queryParams != nil {
		uri += "?" + queryParams.Encode()
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case uri:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		default:
			t.Errorf("Wanted URL path: %v, got instead: %v", uri, r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	ctx, c, teardown := getServerWithRbacAPIClient(t, h, baseURLSuffix)
	defer teardown()

	roles := []string{"LOL"}
	r := struct {
		Roles []string `json:"roles"`
	}{
		roles,
	}
	body, _ := json.Marshal(&r)

	resp, err := request(ctx, c, URLSuffix, queryParams, bytes.NewReader(body))
	if err != nil {
		t.Errorf("Expected no error, got error instead: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status code %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func getServerWithRbacAPIClient(t *testing.T, handlerFunc http.HandlerFunc, baseURLPrefix string) (context.Context, APIClient, func()) {
	s := httptest.NewServer(handlerFunc)

	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fatalf("Invalid test server url: %v", err)
	}
	u.Path = path.Join(u.Path, baseURLPrefix)
	s.URL = u.String()

	c := initRbacAPIClient(t, s.URL)
	ctx := WithAPIKey(context.Background(), "some token")
	return ctx, c, func() {
		s.Close()
	}
}

func initRbacAPIClient(t *testing.T, baseURL string) *Client {
	userAgent := "RbacUserProfileLoader"
	c, err := InitClient(baseURL, userAgent, nil, DefaultTimeout)
	if err != nil {
		t.Errorf("Error initializing client: %v", err)
	}
	return c
}
