package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var envMap = map[string]string{
	port:                defaultPort,
	characterDatabase:   defaultCharacterDatabase,
	characterCollection: defaultCharacterCollection,
	characterArchive:    defaultCharacterArchive,
	logLevel:            defaultlogLevel,
}

//Config is the general struct for app configuration
type Config struct {
	Port                string       `json:"port"`
	CharacterDatabase   string       `json:"characterDatabase"`
	CharacterCollection string       `json:"characterCollection"`
	CharacterArchive    string       `json:"characterArchive"`
	LogLevel            logrus.Level `json:"log-level"`
}

//Accessor is the interface setup for any configuration accessor
type Accessor interface {
	BindEnv(input ...string) error
	IsSet(key string) bool
	GetString(key string) string
}

//New sets up a new config based on the interface passed
func New(accessor Accessor) (c *Config, err error) {
	error := loadEnvVars(accessor)
	if error != nil {
		return nil, error
	}

	currentLogLevel, err := logrus.ParseLevel(envMap[logLevel])
	if err != nil {
		logrus.Warnf("Cannot load log-level: %v", err)
	}

	config := Config{
		Port:                envMap[port],
		CharacterDatabase:   envMap[characterDatabase],
		CharacterCollection: envMap[characterCollection],
		CharacterArchive:    envMap[characterArchive],
		LogLevel:            currentLogLevel,
	}
	return &config, nil
}

func loadEnvVars(accessor Accessor) error {
	for envKey := range envMap {
		err := accessor.BindEnv(envKey)
		if err != nil {
			return fmt.Errorf("error loading environment variable %s: %v", envKey, err)
		}

		if accessor.IsSet(envKey) {
			envMap[envKey] = accessor.GetString(envKey)
		}
	}

	return nil
}
