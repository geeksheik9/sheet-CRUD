package main

import (
	"log"

	"github.com/geeksheik9/sheet-CRUD/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.Info("INITIALIZING SHEET CRUD")

	accessor := viper.New()

	_, err := config.New(accessor)
	if err != nil {
		log.Fatalf("ERROR LOADING CONFIG: %v", err.Error())
	}
}
