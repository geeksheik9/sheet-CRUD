package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	log "github.com/sirupsen/logrus"
)

// ErrorResponse
// struct representation of what gets returned by the RespondWithError function.
// swagger:model
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithError Utility function to convert an error message into a JSON response.
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	// Clean up all quote marks for readability, marshal adds additional "\" escape char in the response JSON.
	RespondWithJSON(w, code, map[string]string{"error": strings.Replace(msg, `"`, ``, -1)})
}

// RespondNoContent Utility function to send a response without any content.
func RespondNoContent(w http.ResponseWriter, code int) {
	if w != nil {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(code)
	}
}

// RespondWithJSON Utility function to convert the payload into a JSON response.
// ORIGINAL:
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	if w != nil {
		response, err := json.Marshal(payload)
		if err != nil {
			log.Errorf("Error in RespondWithJSON marshal: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}
}

// GetJSONRequestBody function to return the request body as string
func GetJSONRequestBody(r *http.Request) (requestBodyJSON string) {
	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	requestBodyJSON = string(bodyBytes)

	return requestBodyJSON
}

// StringToObjectID takes a string and checks if it is a valid objectId hex if so it returns an objectID
func StringToObjectID(ID string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	if objID.IsZero() {
		return objID, errors.New("StringToObjectID() in api.go failed, object ID returned as zero and is not valid")
	}

	return objID, nil
}

// CheckError checks the err message and returns a code based on the message.
func CheckError(err error) int {
	var code int
	// TODO: Check error returns and make this into a switch statement.
	if err == nil {
		code = http.StatusOK
	} else if strings.Contains(err.Error(), "no documents in result") ||
		strings.Contains(err.Error(), "out of bounds") ||
		strings.Contains(err.Error(), "not found") {
		code = http.StatusNotFound
	} else if strings.Contains(err.Error(), "E11000 duplicate key error") ||
		strings.Contains(err.Error(), "E11001 duplicate key error") {
		code = http.StatusConflict
	} else if strings.Contains(err.Error(), "E10334") ||
		strings.Contains(err.Error(), "Invalid request payload, unable to marshal into json, err: ") {
		code = http.StatusBadRequest
	} else {
		code = http.StatusInternalServerError
	}

	return code
}

//BuildQuery sets up the mongo query
func BuildQuery(ID *primitive.ObjectID, name *string, other ...bson.M) bson.M {
	conditions := []bson.M{}
	c := bson.M{}

	if ID != nil {
		conditions = append(conditions, bson.M{"_id": ID})
	}

	if name != nil {
		conditions = append(conditions, bson.M{"name": &name})
	}

	if len(other) > 0 {
		for _, otherFilter := range other {
			conditions = append(conditions, otherFilter)
		}
	}

	if len(conditions) != 0 {
		c = bson.M{"$and": conditions}
	}
	return c
}

//BuildFilter sets up the mongo filtering
func BuildFilter(queryParams url.Values) (int, int, string, bson.M) {
	filter := bson.M{}
	filters := []bson.M{}
	pageNumber := 0
	pageCount := 10000
	sort := "priority"
	if len(queryParams) > 0 {
		for queryParam, paramValue := range queryParams {
			switch queryParam {
			case "pageNumber":
				pageNumber, _ = strconv.Atoi(paramValue[0])
			case "pageCount":
				pageCount, _ = strconv.Atoi(paramValue[0])
			case "sort":
				sort = paramValue[0]
			default:
				m := bson.M{queryParam: paramValue[0]}
				filters = append(filters, m)
			}
		}
		if len(filters) > 0 {
			filter = bson.M{"$and": filters}
		} else {
			filter = nil
		}
		return pageNumber, pageCount, sort, filter
	}

	return pageNumber, pageCount, sort, nil
}
