package db

import (
	"context"
	"errors"
	"net/url"
	"time"

	model "github.com/geeksheik9/sheet-CRUD/models"
	"github.com/geeksheik9/sheet-CRUD/pkg/api"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

//GetForceCharacterSheets returns all FFG Star Wars Force sensitive character sheets that reside in the database
func (d *CharacterDB) GetForceCharacterSheets(queryParams url.Values) ([]model.ForceCharacterSheet, error) {
	logrus.Debugf("BEGIN - GetForceCharacterSheets")

	collection := d.client.Database(d.databaseName).Collection(d.collectionName)

	pageNumber, pageCount, sort, filter := api.BuildFilter(queryParams)
	skip := 0
	if pageNumber > 0 {
		skip = (pageNumber - 1) * pageCount
	}

	opts := options.Find().
		SetMaxTime(30 * time.Second).
		SetSkip(int64(skip)).
		SetLimit(int64(pageCount)).
		SetSort(bson.D{{
			Key:   sort,
			Value: 1,
		}})

	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	matches := []model.ForceCharacterSheet{}

	for cur.Next(context.Background()) {
		elem := model.ForceCharacterSheet{}
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		matches = append(matches, elem)
	}

	return matches, nil
}

//FindForceCharacterSheetByID finds a specific force character sheet by a provided ID
func (d *CharacterDB) FindForceCharacterSheetByID(mongoID primitive.ObjectID) (*model.ForceCharacterSheet, error) {
	logrus.Debug("BEGIN - FindForceCharacterSheet: %v", mongoID)

	collection := d.client.Database(d.databaseName).Collection(d.collectionName)
	query := api.BuildQuery(&mongoID, nil)

	sheet := model.ForceCharacterSheet{}

	err := collection.FindOne(context.Background(), query).Decode(&sheet)
	if err != nil {
		return nil, err
	}

	return &sheet, err
}

//UpdateForceCharacterSheetByID updates a specific force character sheet by provided ID
func (d *CharacterDB) UpdateForceCharacterSheetByID(sheet model.ForceCharacterSheet, mongoID primitive.ObjectID) error {
	logrus.Debug("BEGIN - UpdateForceCharacterSheetByID: %v", mongoID)

	collection := d.client.Database(d.databaseName).Collection(d.collectionName)

	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": mongoID}, bson.D{{
		Key:   "$set",
		Value: sheet,
	}})
	if err != nil {
		return err
	}

	logrus.Infof("result: %v", result)

	if result.MatchedCount != 1 {
		return errors.New("Could not update sheet tried to update " + mongoID.Hex() + " got " + string(result.MatchedCount) + " matches instead of 1")
	}

	if result.ModifiedCount != 1 {
		return errors.New("Could not update sheet tried to updated " + mongoID.Hex() + " tried to update " + string(result.ModifiedCount) + " number of results instead of 1")
	}

	return nil
}
