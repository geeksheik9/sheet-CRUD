package rbac

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	model "github.com/geeksheik9/sheet-CRUD/models"
	"github.com/geeksheik9/sheet-CRUD/pkg/api"
	"github.com/sirupsen/logrus"
)

// RoleAPIClient is an interface for the rbac call
type RoleAPIClient interface {
	PerformRBACCheck(ctx context.Context, tokenString string, roles []model.Role) (bool, error)
}

// RoleClient is the implementation for the RoleAPIClient
type RoleClient struct {
	Client api.APIClient
}

//NewRoleClient creates a new instance of the role client
func NewRoleClient(baseURL string, userAgent string, timeout time.Duration, httpClient *http.Client) (*RoleClient, error) {
	client, err := api.InitClient(baseURL, userAgent, httpClient, timeout)
	return &RoleClient{Client: client}, err
}

// PerformRBACCheck checks if a user is authorized to make a request
func (c *RoleClient) PerformRBACCheck(ctx context.Context, tokenString string, roles []model.Role) (bool, error) {
	logrus.Debug("BEGIN - PerformRBACCheck")
	body, _ := json.Marshal(roles)

	result := api.Handler(c.Client.Post(api.WithAPIKey(ctx, tokenString), "/check-role", bytes.NewReader(body)))
	if result.Code == 404 {
		return false, result.Error
	}

	var user string
	err := json.Unmarshal(result.Raw, &user)
	if err != nil {
		return false, nil
	}

	authorized := true
	if user == "User does not have the required roles" {
		authorized = false
	}

	return authorized, nil
}
