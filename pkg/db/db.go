package db

import (
	"context"

	model "github.com/geeksheik9/sheet-CRUD/models"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//CharacterDB is the data access object for the Star Wars FFG character sheets
type CharacterDB struct {
	client         *mongo.Client
	databaseName   string
	collectionName string
	archiveName    string
}

//Ping checks that the database is running
func (d *CharacterDB) Ping() error {
	err := d.client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		logrus.Errorf("ERROR connectiong to database %v", err)
	}
	return err
}

//InsertForceCharacterSheet inserts the FFG Star Wars Force sensitive character sheet into the database
func (d *CharacterDB) InsertForceCharacterSheet(sheet model.ForceCharacterSheet) error {
	logrus.Debug("BEGIN - InsertForceCharacterSheet")

	collection := d.client.Database(d.databaseName).Collection(d.collectionName)

	_, err := collection.InsertOne(context.Background(), sheet)

	return err
}
