package handler

import (
	"net/http"

	model "github.com/geeksheik9/sheet-CRUD/models"
	"github.com/geeksheik9/sheet-CRUD/pkg/api"
	"github.com/gorilla/mux"
)

//CharacterDatabase is the interface setup for accesssing the character database
type CharacterDatabase interface {
	GetCharacterSheets()
	FindCharacterSheetByID()
	FindCharacterSheetByName()
	UpdateCharacterSheetByID()
	InsertCharacterSheet()
	DeleteCharacterSheetByID()
	FindArchivedCharacterSheetByID()
	FindArchivedCharacterSheetByName()
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

		response := model.HealthCheckResponse{
			APIVersion: s.Version,
			DBError:    dbErr.Error(),
		}

		if dbErr != nil {
			api.RespondWithJSON(w, http.StatusFailedDependency, response)
			return
		}

		api.RespondWithJSON(w, http.StatusOK, response)
	})
}
