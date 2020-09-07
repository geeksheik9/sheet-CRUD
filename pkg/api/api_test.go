package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_RespondWithJson(t *testing.T) {
	w := httptest.NewRecorder()
	RespondWithJSON(w, http.StatusOK, map[string]string{"error": "Test payload"})
	if w.Code != http.StatusOK {
		t.Errorf("TestRespondWithJson() error:\n   expected: %v\n   got:      %d", http.StatusOK, w.Code)
	}
}

func Test_RespondWithError(t *testing.T) {
	w := httptest.NewRecorder()
	RespondWithError(w, http.StatusNotFound, "Test payload")
	if w.Code != http.StatusNotFound {
		t.Errorf("TestRespondWithError() error:\n   expected: %v\n   got:      %d", http.StatusNotFound, w.Code)
	}
}

func Test_RespondWithError_RemoveQuotes(t *testing.T) {
	w := httptest.NewRecorder()
	RespondWithError(w, http.StatusNotFound, `Te"st ""payl"o"ad`)
	body := w.Body.String()
	if body != `{"error":"Test payload"}` {
		t.Errorf("RespondNoContent():\n   expected error body:\n   got:      %s", body)
	}
}

func Test_RespondWithError_Empty(t *testing.T) {
	w := httptest.NewRecorder()
	RespondWithError(w, http.StatusNotFound, "")
	body := w.Body.String()
	if body != `{"error":""}` {
		t.Errorf("RespondNoContent() error:\n   expected empty error body:\n   got:      %s", body)
	}
}

func Test_RespondNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	RespondNoContent(w, http.StatusOK)
	if w.Code != http.StatusOK {
		t.Errorf("TestRespondWithJson() error:\n   expected: %v\n   got:      %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "" {
		t.Errorf("RespondNoContent() error:\n   expected empty body:\n   got:      %s", w.Body.String())
	}
}

func Test_GetJSONRequestBody(t *testing.T) {
	s := `{"hello":"goodbye"}`
	r, _ := http.NewRequest("GET", "/any", bytes.NewBuffer([]byte(s)))
	res := GetJSONRequestBody(r)

	if s != res {
		t.Errorf("GetJSONRequestBody() error:\n   expected:%s\n   got:      %s", s, res)
	}
}

func Test_GetJSONRequestBody_Malformed(t *testing.T) {
	s := `{"hello":"goodbye"`
	r, _ := http.NewRequest("GET", "/any", bytes.NewBuffer([]byte(s)))
	res := GetJSONRequestBody(r)

	if s != res {
		t.Errorf("GetJSONRequestBody() error:\n   expected:%s\n   got:      %s", s, res)
	}
}

func Test_StringToObjectID_notID(t *testing.T) {
	testString := "asdf"

	expectedErr := fmt.Errorf("%v is not a valid objectID", testString)

	_, err := StringToObjectID(testString)
	if err == nil {
		t.Errorf("StringToObjectID() error:\n   expected: %v\n  got:    %v", expectedErr, err)
	}
}

func Test_StringToObjectID_isID(t *testing.T) {
	testString := "5b883e25ad3d111aa02b4693"

	var expectedErr error

	_, err := StringToObjectID(testString)
	if err != nil {
		t.Errorf("StringToObjectID() error:\n   expected: %v\n  got:    %v", expectedErr, err)
	}
}

func TestCheckError(t *testing.T) {
	if code := CheckError(nil); code != http.StatusOK {
		t.Errorf("TestCheckError(),\n   expected: %v\n   got:      %v", http.StatusOK, code)
	}
	if code := CheckError(errors.New("no documents in result")); code != http.StatusNotFound {
		t.Errorf("TestCheckError(),\n   expected: %v\n   got:      %v", http.StatusNotFound, code)
	}
	if code := CheckError(errors.New("E11000 duplicate key error")); code != http.StatusConflict {
		t.Errorf("TestCheckError(),\n   expected: %v\n   got:      %v", http.StatusConflict, code)
	}
	if code := CheckError(errors.New("E10334")); code != http.StatusBadRequest {
		t.Errorf("TestCheckError(),\n   expected: %v\n   got:      %v", http.StatusBadRequest, code)
	}
	if code := CheckError(errors.New("E1")); code != http.StatusInternalServerError {
		t.Errorf("TestCheckError(),\n   expected: %v\n   got:      %v", http.StatusInternalServerError, code)
	}
}
