package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	model "github.com/geeksheik9/sheet-CRUD/models"
	"github.com/geeksheik9/sheet-CRUD/pkg/db/mocks"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InitMockCharacterService(sheetsToReturn []model.ForceCharacterSheet, sheetToReturn *model.ForceCharacterSheet, errorToReturn error) (s CharacterService) {
	db := mocks.MockCharacterDB{
		SheetsToReturn: sheetsToReturn,
		SheetToReturn:  sheetToReturn,
		ErrorToReturn:  errorToReturn,
	}

	return CharacterService{
		Version:  "test",
		Database: &db,
	}
}

func mockCharacter(id primitive.ObjectID, name string, threshold int64, current int64, rating int64) model.ForceCharacterSheet {
	character := model.ForceCharacterSheet{
		ID:            id,
		CharacterName: name,
		PlayerName:    name,
		Wounds: model.Amount{
			Threshold: threshold,
			Current:   current,
		},
		ForceRating: rating,
	}
	return character
}

func mockCharacters(character model.ForceCharacterSheet) []model.ForceCharacterSheet {
	characters := []model.ForceCharacterSheet{
		character,
	}
	return characters
}

func TestCharacterService_GetForceCharacterSheets_Success(t *testing.T) {
	id := primitive.NewObjectID()
	sheet := mockCharacter(id, "test", 2, 0, 5)
	sheets := mockCharacters(sheet)
	service := InitMockCharacterService(sheets, nil, nil)

	r, err := http.NewRequest("GET", "/force-character-sheet", nil)
	if err != nil {
		t.Errorf("GetForceCharacterSheets() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("GetForceCharacterSheets() error:\ngot: %v\nexpected: %v", w.Code, http.StatusOK)
	}

	resp := []model.ForceCharacterSheet{}
	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Errorf("GetForceCharacterSheets() json decode error:\ngot: %v\nexpected: <nil>", err)
	}
	if resp[0].ID != sheet.ID || resp[0].CharacterName != sheet.CharacterName || resp[0].PlayerName != sheet.PlayerName ||
		resp[0].Wounds.Threshold != sheet.Wounds.Threshold || resp[0].Wounds.Current != sheet.Wounds.Current || resp[0].ForceRating != sheet.ForceRating {
		t.Errorf("GetForceCharacterSheets() error:\n got:%v\nexpected:%v", resp[0], sheet)
	}
}

func TestCharacterService_GetForceCharacterSheets_DBError(t *testing.T) {
	service := InitMockCharacterService(nil, nil, errors.New("test error"))

	r, err := http.NewRequest("GET", "/force-character-sheet", nil)
	if err != nil {
		t.Errorf("GetForceCharacterSheets() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetForceCharacterSheets() error:\ngot: %v\nexpected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestCharacterService_FindForceCharacterSheetByID_Success(t *testing.T) {
	id := primitive.NewObjectID()
	sheet := mockCharacter(id, "test", 2, 0, 5)
	service := InitMockCharacterService(nil, &sheet, nil)

	r, err := http.NewRequest("GET", "/force-character-sheet/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("FindForceCharacterSheetByID() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("FindForceCharacterSheetByID() error:\ngot:%v\nexpected:%v", w.Code, http.StatusOK)
	}

	resp := model.ForceCharacterSheet{}
	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Errorf("FindForceCharacterSheetByID() json decode error:\ngot: %v\nexpected: <nil>", err)
	}
	if resp.ID != sheet.ID || resp.CharacterName != sheet.CharacterName || resp.PlayerName != sheet.PlayerName ||
		resp.Wounds.Threshold != sheet.Wounds.Threshold || resp.Wounds.Current != sheet.Wounds.Current || resp.ForceRating != sheet.ForceRating {
		t.Errorf("FindForceCharacterSheetByID() error:\n got:%v\nexpected:%v", resp, sheet)
	}
}

func TestCharacterService_FindForceCharacterSheetByID_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockCharacterService(nil, nil, errors.New("test error"))

	r, err := http.NewRequest("GET", "/force-character-sheet/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("FindForceCharacterSheetByID() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("FindForceCharacterSheetByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestCharacterService_FindForceCharacterSheetByID_BadID(t *testing.T) {
	id := "this is a bad id"
	service := InitMockCharacterService(nil, nil, nil)

	r, err := http.NewRequest("GET", "/force-character-sheet/"+id, nil)
	if err != nil {
		t.Errorf("FindForceCharacterSheetByID() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("FindForceCharacterSheetByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestCharacterService_UpdateForceCharacterSheetByID_Success(t *testing.T) {
	id := primitive.NewObjectID()
	sheet := mockCharacter(id, "test", 2, 0, 5)
	service := InitMockCharacterService(nil, &sheet, nil)

	request, _ := json.Marshal(sheet)

	r, err := http.NewRequest("PUT", "/force-character-sheet/"+id.Hex(), bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateForceCharacterSheetByID() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("UpdateForceCharacterSheetByID() error:\ngot:%v\nexpected:%v", w.Code, http.StatusOK)
	}
}

func TestCharacterService_UpdateForceCharacterSheetByID_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	sheet := mockCharacter(id, "test", 2, 0, 5)
	service := InitMockCharacterService(nil, nil, errors.New("test error"))

	request, _ := json.Marshal(sheet)

	r, err := http.NewRequest("PUT", "/force-character-sheet/"+id.Hex(), bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateForceCharacterSheetByID() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("UpdateForceCharacterSheetByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestCharacterService_UpdateForceCharacterSheetByID_BadID(t *testing.T) {
	id := "this is a bad id"
	sheet := mockCharacter(primitive.NewObjectID(), "test", 2, 0, 5)
	service := InitMockCharacterService(nil, nil, nil)

	request, _ := json.Marshal(sheet)

	r, err := http.NewRequest("PUT", "/force-character-sheet/"+id, bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("FindForceCharacterSheetByID() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("FindForceCharacterSheetByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestCharacterService_UpdateForceCharacterSheetByID_BadJSON(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockCharacterService(nil, nil, nil)

	request, _ := json.Marshal(`{bad json`)

	r, err := http.NewRequest("PUT", "/force-character-sheet/"+id.Hex(), bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateForceCharacterSheetByID() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("UpdateForceCharacterSheetByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusBadRequest)
	}
}

func TestCharacterService_InsertForceCharacterSheet_Success(t *testing.T) {
	id := primitive.NewObjectID()
	sheet := mockCharacter(id, "test", 2, 0, 5)
	service := InitMockCharacterService(nil, &sheet, nil)

	request, _ := json.Marshal(sheet)

	r, err := http.NewRequest("POST", "/force-character-sheet", bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("InsertForceCharacterSheet() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Errorf("InsertForceCharacterSheet() error:\ngot:%v\nexpected:%v", w.Code, http.StatusOK)
	}
}

func TestCharacterService_InsertForceCharacterSheet_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	sheet := mockCharacter(id, "test", 2, 0, 5)
	service := InitMockCharacterService(nil, nil, errors.New("test error"))

	request, _ := json.Marshal(sheet)

	r, err := http.NewRequest("POST", "/force-character-sheet", bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("InsertForceCharacterSheet() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("InsertForceCharacterSheet() error:\ngot: %v\nexpected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestCharacterService_InsertForceCharacterSheet_BadJSON(t *testing.T) {
	service := InitMockCharacterService(nil, nil, nil)

	request, _ := json.Marshal(`{bad json`)

	r, err := http.NewRequest("POST", "/force-character-sheet", bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateForceCharacterSheetByID() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("UpdateForceCharacterSheetByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusBadRequest)
	}
}

func TestCharacterService_DeleteForceCharacterSheetByID_Success(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockCharacterService(nil, nil, nil)

	r, err := http.NewRequest("DELETE", "/force-character-sheet/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("DeleteForceCharacterSheetByID() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusNoContent {
		t.Errorf("DeleteForceCharacterSheetByID() error:\ngot:%v\nexpected:%v", w.Code, http.StatusOK)
	}
}

func TestCharacterService_DeleteForceCharacterSheetByID_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockCharacterService(nil, nil, errors.New("test error"))

	r, err := http.NewRequest("DELETE", "/force-character-sheet/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("DeleteForceCharacterSheetByID() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("DeleteForceCharacterSheetByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestCharacterService_DeleteForceCharacterSheetByID_BadID(t *testing.T) {
	id := "this is a bad id"
	service := InitMockCharacterService(nil, nil, nil)

	r, err := http.NewRequest("DELETE", "/force-character-sheet/"+id, nil)
	if err != nil {
		t.Errorf("FindForceCharacterSheetByID() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("FindForceCharacterSheetByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestCharacterService_PingCheck(t *testing.T) {
	service := InitMockCharacterService(nil, nil, nil)

	r, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Errorf("Ping() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Ping() error:\ngot: %v\n expected: %v", w.Code, http.StatusOK)
	}
}

func TestCharacterService_HealthCheck(t *testing.T) {
	service := InitMockCharacterService(nil, nil, nil)
	r, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Errorf("Ping() error creating request:\ngot: %v\nexpected:<no error>", err)
	}
	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Ping() error:\ngot: %v\n expected: %v", w.Code, http.StatusOK)
	}
}

func TestCharacterService_HealthCheck_error(t *testing.T) {
	service := InitMockCharacterService(nil, nil, errors.New("test error"))
	r, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Errorf("Ping() error creating request:\ngot: %v\nexpected:<no error>", err)
	}
	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusFailedDependency {
		t.Errorf("Ping() error:\ngot: %v\n expected: %v", w.Code, http.StatusFailedDependency)
	}
}
