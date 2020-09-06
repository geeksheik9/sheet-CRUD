package config

import (
	"errors"
	"testing"

	"github.com/geeksheik9/sheet-CRUD/config/mocks"
)

func TestConfig_NewNoError(t *testing.T) {
	configAccessor := &mocks.ConfigAccessor{}

	configAccessor.On("BindEnv", port).Return(nil)
	configAccessor.On("IsSet", port).Return(false)

	for envKey := range envMap {
		if envKey != port {
			configAccessor.On("BindEnv", envKey).Return(nil)
			configAccessor.On("IsSet", envKey).Return(true)
			configAccessor.On("GetString", envKey).Return("dummyEnvValue")
		}
	}

	c, _ := New(configAccessor)

	configAccessor.AssertNumberOfCalls(t, "BindEnv", len(envMap))
	configAccessor.AssertNumberOfCalls(t, "IsSet", len(envMap))
	configAccessor.AssertNumberOfCalls(t, "GetString", len(envMap)-1)

	// Test that port uses default value
	if c.Port != defaultPort {
		t.Errorf("Environment variable PORT returned wrong value: got %v, want %v", c.Port, defaultPort)
	}
}

func TestConfig_NewWithError(t *testing.T) {
	configAccessor := &mocks.ConfigAccessor{}

	for envKey := range envMap {
		if envKey == characterDatabase {
			configAccessor.On("BindEnv", envKey).Return(errors.New("test error"))
		} else {
			configAccessor.On("BindEnv", envKey).Return(nil)
			configAccessor.On("IsSet", envKey).Return(false)
		}
	}

	expectedErr := "error loading environment variable CHARACTER_DATABASE: test error"
	_, err := New(configAccessor)
	if err.Error() != expectedErr {
		t.Errorf("New() returned wrong value: got %v, want %v", err, expectedErr)
	}
}
