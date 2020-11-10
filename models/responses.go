package model

//HealthCheckResponse returns the version of the api and whether or not the DB is queryable
type HealthCheckResponse struct {
	APIVersion string `json:"apiVersion"`
	DBError    string `json:"dbError"`
}

//APIResponse is the basic response from the APIClient
type APIResponse struct {
	Body       []byte
	StatusCode int
}

// Role is the implementation of roles that a user would have
type Role struct {
	Name string `json:"name" bson:"name"`
}
