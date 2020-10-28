package mocks

import (
	"net/url"

	model "github.com/geeksheik9/sheet-CRUD/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//MockCharacterDB is the mock struct for testing
type MockCharacterDB struct {
	SheetsToReturn []model.ForceCharacterSheet
	SheetToReturn  *model.ForceCharacterSheet
	ErrorToReturn  error
}

//GetForceCharacterSheets is the mock implementation for testing
func (db *MockCharacterDB) GetForceCharacterSheets(query url.Values) ([]model.ForceCharacterSheet, error) {
	return db.SheetsToReturn, db.ErrorToReturn
}

//FindForceCharacterSheetByID is the mock implementation for testing
func (db *MockCharacterDB) FindForceCharacterSheetByID(mongoID primitive.ObjectID) (*model.ForceCharacterSheet, error) {
	return db.SheetToReturn, db.ErrorToReturn
}

//UpdateForceCharacterSheetByID is the mock implementation for testing
func (db *MockCharacterDB) UpdateForceCharacterSheetByID(sheet model.ForceCharacterSheet, mongoID primitive.ObjectID) error {
	return db.ErrorToReturn
}

//InsertForceCharacterSheet is the mock implementation for testing
func (db *MockCharacterDB) InsertForceCharacterSheet(sheet model.ForceCharacterSheet) error {
	return db.ErrorToReturn
}

//DeleteForceCharacterSheetByID is the mock implementation for testing
func (db *MockCharacterDB) DeleteForceCharacterSheetByID(mongoID primitive.ObjectID) error {
	return db.ErrorToReturn
}

//Ping is the mock implementation for testing
func (db *MockCharacterDB) Ping() error {
	return db.ErrorToReturn
}
