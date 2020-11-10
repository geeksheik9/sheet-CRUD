package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	model "github.com/geeksheik9/sheet-CRUD/models"
	"github.com/geeksheik9/sheet-CRUD/pkg/api"
	"github.com/geeksheik9/sheet-CRUD/pkg/rbac"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CharacterDatabase is the interface setup for accesssing the character database
type CharacterDatabase interface {
	GetForceCharacterSheets(query url.Values) ([]model.ForceCharacterSheet, error)
	FindForceCharacterSheetByID(mongoID primitive.ObjectID) (*model.ForceCharacterSheet, error)
	UpdateForceCharacterSheetByID(sheet model.ForceCharacterSheet, mongoID primitive.ObjectID) error
	InsertForceCharacterSheet(sheet model.ForceCharacterSheet) error
	DeleteForceCharacterSheetByID(mongoID primitive.ObjectID) error
	Ping() error
}

//CharacterService is the implementation of the service to access character sheets
type CharacterService struct {
	Version    string
	Database   CharacterDatabase
	RBACClient rbac.RoleAPIClient
}

//Routes sets up the routes for the RESTful interface
func (s *CharacterService) Routes(r *mux.Router) *mux.Router {
	r.HandleFunc("/ping", s.PingCheck).Methods(http.MethodGet)
	r.Handle("/health", s.healthCheck(s.Database)).Methods(http.MethodGet)

	// swagger:route POST /force-character-sheet ForceCharacterSheet
	//
	// Insert Force Character Sheet
	//
	// Consumes:
	// - application/json
	// Schemes: http, https
	//
	// responses:
	// 201: description:Created
	// 400: description:Bad request
	// 500: description:Internal Server Error
	r.HandleFunc("/force-character-sheet", s.InsertForceCharacterSheet).Methods(http.MethodPost)
	// swagger:route GET /force-character-sheet ForceCharacterSheet
	//
	// Get Force Character Sheet
	//
	// Consumes:
	// - application/json
	// Schemes: http, https
	//
	// responses:
	// 200: []ForceCharacterSheet
	// 400: description:Bad request
	// 404: description:No records
	// 500: description:Internal Server Error
	r.HandleFunc("/force-character-sheet", s.GetForceCharacterSheets).Methods(http.MethodGet)
	// swagger:route GET /force-character-sheet ForceCharacterSheet
	//
	// Get Force Character Sheet by ID
	//
	// Consumes:
	// - application/json
	// Schemes: http, https
	//
	// responses:
	// 200: ForceCharacterSheet
	// 400: description:Bad request
	// 404: description:No records
	// 500: description:Internal Server Error
	r.HandleFunc("/force-character-sheet/{ID}", s.FindForceCharacterSheetByID).Methods(http.MethodGet)
	// swagger:route PUT /force-character-sheet/{ID} ForceCharacterSheet
	//
	// Update Force Character Sheet by ID
	//
	// Consumes:
	// - application/json
	// Schemes: http, https
	//
	// responses:
	// 200: description:Success
	// 400: description:Bad request
	// 404: description:No records
	// 500: description:Internal Server Error
	r.HandleFunc("/force-character-sheet/{ID}", s.UpdateForceCharacterSheetByID).Methods(http.MethodPut)
	// swagger:route DELETE /force-character-sheet/{ID} ForceCharacterSheet
	//
	// Update Force Character Sheet by ID
	//
	// Consumes:
	// - application/json
	// Schemes: http, https
	//
	// responses:
	// 204: description:No Content
	// 400: description:Bad request
	// 404: description:No records
	// 500: description:Internal Server Error
	r.HandleFunc("/force-character-sheet/{ID}", s.DeleteForceCharacterSheetByID).Methods(http.MethodDelete)

	fs := http.FileServer(http.Dir("./swagger-ui/"))
	r.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger", fs))

	return r
}

//PingCheck checks that the app is running and returns 200, OK, version
func (s *CharacterService) PingCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK, " + s.Version))
}

func (s *CharacterService) healthCheck(database CharacterDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbErr := database.Ping()
		var stringDBErr string

		if dbErr != nil {
			stringDBErr = dbErr.Error()
		}

		response := model.HealthCheckResponse{
			APIVersion: s.Version,
			DBError:    stringDBErr,
		}

		if dbErr != nil {
			api.RespondWithJSON(w, http.StatusFailedDependency, response)
			return
		}

		api.RespondWithJSON(w, http.StatusOK, response)
	})
}

//InsertForceCharacterSheet is the handler function for inserting a force character sheet
func (s *CharacterService) InsertForceCharacterSheet(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("InsertForceCharacterSheet invoked with url: %v", r.URL)
	defer r.Body.Close()

	tokenString := r.Header.Get("Authorization")
	if strings.Contains(tokenString, "Bearer") {
		tokenString = strings.Trim(tokenString, "Bearer")
		tokenString = strings.Trim(tokenString, " ")
	}
	if tokenString == "" {
		api.RespondWithError(w, http.StatusUnauthorized, "User is not authorized to make this request")
		return
	}

	roles := []model.Role{
		{Name: "gamemaster"},
	}
	authorized, err := s.RBACClient.PerformRBACCheck(context.Background(), tokenString, roles)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}
	if !authorized {
		api.RespondWithJSON(w, http.StatusUnauthorized, "User is not authorized to access this resource")
		return
	}

	var characterSheet model.ForceCharacterSheet
	characterSheet.ID = primitive.NewObjectID()

	err = json.NewDecoder(r.Body).Decode(&characterSheet)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	if characterSheet.Version == 0 {
		characterSheet.Version = 1
	}

	err = s.Database.InsertForceCharacterSheet(characterSheet)
	if err != nil {
		api.RespondWithJSON(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, characterSheet.ID)
}

//GetForceCharacterSheets is the handler function for getting all force character sheet
func (s *CharacterService) GetForceCharacterSheets(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("GetForceCharacterSheet invoked with url: %v", r.URL)

	tokenString := r.Header.Get("Authorization")
	if strings.Contains(tokenString, "Bearer") {
		tokenString = strings.Trim(tokenString, "Bearer")
		tokenString = strings.Trim(tokenString, " ")
	}
	if tokenString == "" {
		api.RespondWithError(w, http.StatusUnauthorized, "User is not authorized to make this request")
		return
	}

	roles := []model.Role{
		{Name: "viewer"},
	}
	authorized, err := s.RBACClient.PerformRBACCheck(context.Background(), tokenString, roles)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}
	if !authorized {
		api.RespondWithJSON(w, http.StatusUnauthorized, "User is not authorized to access this resource")
		return
	}

	sheets, err := s.Database.GetForceCharacterSheets(r.URL.Query())
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, sheets)
}

//FindForceCharacterSheetByID is the handler function for getting a specific character sheet by database ID
func (s *CharacterService) FindForceCharacterSheetByID(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("BEGIN - FindCharacterSheetByID invoked with url: %v", r.URL)

	tokenString := r.Header.Get("Authorization")
	if strings.Contains(tokenString, "Bearer") {
		tokenString = strings.Trim(tokenString, "Bearer")
		tokenString = strings.Trim(tokenString, " ")
	}
	if tokenString == "" {
		api.RespondWithError(w, http.StatusUnauthorized, "User is not authorized to make this request")
		return
	}

	roles := []model.Role{
		{Name: "viewer"},
	}
	authorized, err := s.RBACClient.PerformRBACCheck(context.Background(), tokenString, roles)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}
	if !authorized {
		api.RespondWithJSON(w, http.StatusUnauthorized, "User is not authorized to access this resource")
		return
	}

	vars := mux.Vars(r)
	ID := vars["ID"]

	objectID, err := api.StringToObjectID(ID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	sheet, err := s.Database.FindForceCharacterSheetByID(objectID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, sheet)
}

//UpdateForceCharacterSheetByID is the handler function for updating a specific character sheet by database ID
func (s *CharacterService) UpdateForceCharacterSheetByID(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("BEGIN - UpdateCharacterSheetByID invoked with url: %v", r.URL)
	defer r.Body.Close()

	vars := mux.Vars(r)
	ID := vars["ID"]

	tokenString := r.Header.Get("Authorization")
	if strings.Contains(tokenString, "Bearer") {
		tokenString = strings.Trim(tokenString, "Bearer")
		tokenString = strings.Trim(tokenString, " ")
	}
	if tokenString == "" {
		api.RespondWithError(w, http.StatusUnauthorized, "User is not authorized to make this request")
		return
	}

	roles := []model.Role{
		{Name: "player"},
	}
	authorized, err := s.RBACClient.PerformRBACCheck(context.Background(), tokenString, roles)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}
	if !authorized {
		api.RespondWithJSON(w, http.StatusUnauthorized, "User is not authorized to access this resource")
		return
	}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	sheet := model.ForceCharacterSheet{}
	err = json.NewDecoder(r.Body).Decode(&sheet)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.Database.UpdateForceCharacterSheetByID(sheet, objectID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, objectID)
}

//DeleteForceCharacterSheetByID is the handler function for deleting a specific character sheet by database ID
func (s *CharacterService) DeleteForceCharacterSheetByID(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("BEGIN - UpdateCharacterSheetByID invoked with url: %v", r.URL)

	vars := mux.Vars(r)
	ID := vars["ID"]

	tokenString := r.Header.Get("Authorization")
	if strings.Contains(tokenString, "Bearer") {
		tokenString = strings.Trim(tokenString, "Bearer")
		tokenString = strings.Trim(tokenString, " ")
	}
	if tokenString == "" {
		api.RespondWithError(w, http.StatusUnauthorized, "User is not authorized to make this request")
		return
	}

	roles := []model.Role{
		{Name: "gamemaster"},
	}
	authorized, err := s.RBACClient.PerformRBACCheck(context.Background(), tokenString, roles)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}
	if !authorized {
		api.RespondWithJSON(w, http.StatusUnauthorized, "User is not authorized to access this resource")
		return
	}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	err = s.Database.DeleteForceCharacterSheetByID(objectID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
	}

	api.RespondNoContent(w, http.StatusNoContent)
}
