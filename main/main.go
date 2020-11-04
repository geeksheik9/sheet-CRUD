// Package main Sheet-CRUD API
//
// API for creating, reading, updating, and deleting FFG star wars character sheets
//
// version: 0.0.2-alpha
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/geeksheik9/sheet-CRUD/config"
	"github.com/geeksheik9/sheet-CRUD/pkg/db"
	"github.com/geeksheik9/sheet-CRUD/pkg/handler"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var version string

func main() {
	//go:generate swagger generate spec
	logrus.Info("INITIALIZING SHEET CRUD")

	accessor := viper.New()

	config, err := config.New(accessor)
	if err != nil {
		log.Fatalf("ERROR LOADING CONFIG: %v", err.Error())
	}

	timeout := time.Second * 5
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := db.InitializeClients(ctx)
	if err != nil {
		logrus.Warnf("Failed to intialize client with error: %v, trying again", err)
		err = nil
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*60)
		client, err = db.InitializeClients(ctx)
		if err != nil {
			logrus.Fatalf("Failed to initialize database client a second time with error: %v", err)
		}
	}

	defer client.Disconnect(context.Background())

	database := db.InitializeDatabases(client, config)
	if database == nil {
		log.Fatalf("Error no database from client %v", client)
	}

	characterService := handler.CharacterService{
		Version:  version,
		Database: database,
	}

	r := mux.NewRouter().StrictSlash(true)

	r = characterService.Routes(r)
	fmt.Printf("Sever listening on port %v\n", config.Port)
	logrus.Info("END")
	log.Fatal(http.ListenAndServe(":"+config.Port, cors.Default().Handler(r)))
}
