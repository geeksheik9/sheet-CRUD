package handler

import (
	"encoding/json"
	"net/http"

	model "github.com/geeksheik9/sheet-CRUD/models"
	"github.com/geeksheik9/sheet-CRUD/pkg/api"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CharacterDatabase is the interface setup for accesssing the character database
type CharacterDatabase interface {
	//GetCharacterSheets()
	//FindCharacterSheetByID()
	//FindCharacterSheetByName()
	//UpdateCharacterSheetByID()
	InsertForceCharacterSheet(model.ForceCharacterSheet) error
	//DeleteCharacterSheetByID()
	//FindArchivedCharacterSheetByID()
	//FindArchivedCharacterSheetByName()
	Ping() error
}

//CharacterService is the implementation of the service to access character sheets
type CharacterService struct {
	Version  string
	Database CharacterDatabase
}

//Routes sets up the routes for the RESTful interface
func (s *CharacterService) Routes(r *mux.Router) *mux.Router {
	r.HandleFunc("/ping", s.PingCheck).Methods(http.MethodGet)
	r.Handle("/health", s.healthCheck(s.Database)).Methods(http.MethodGet)
	r.HandleFunc("/force-character-sheet", s.InsertForceCharacterSheet).Methods(http.MethodPost)

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

		logrus.Infof("Log DB err: %v", dbErr)
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

//InsertForceCharacterSheet is the handler function for inserting a character sheet
func (s *CharacterService) InsertForceCharacterSheet(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("InsertForceCharacterSheet invoked with url: %v", r.URL)

	var characterSheet model.ForceCharacterSheet
	characterSheet.ID = primitive.NewObjectID()

	err := json.NewDecoder(r.Body).Decode(&characterSheet)
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
