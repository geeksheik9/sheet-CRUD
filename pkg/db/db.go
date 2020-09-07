package db

import (
	"context"

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
