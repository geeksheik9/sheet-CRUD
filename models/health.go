package model

//HealthCheckResponse returns the version of the api and whether or not the DB is queryable
type HealthCheckResponse struct {
	APIVersion string `json:"apiVersion"`
	DBError    string `json:"dbError"`
}
