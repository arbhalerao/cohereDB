package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// Note(aditya): Temporarily setting all environment variables through code until Docker Compose is implemented
func setEnvs() error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	envVars := map[string]string{
		"HTTP_ADDR":      "0.0.0.0:8080",
		"BADGER_DB_PATH": filepath.Join(pwd, "../db/db.db"),
	}

	for key, value := range envVars {
		err := os.Setenv(key, value)
		if err != nil {
			return fmt.Errorf("error setting environment variable: %v", err)
		}
	}

	return nil
}

func LoadConfig(configs ...interface{}) error {
	// Note(aditya): Temporarily setting all environment variables through code until Docker Compose is implemented
	if err := setEnvs(); err != nil {
		fmt.Println(err)
	}

	for _, config := range configs {
		v := reflect.ValueOf(config).Elem()
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			envVar := field.Tag.Get("env")
			if envVar == "" {
				continue
			}
			value := os.Getenv(envVar)
			if value == "" {
				return fmt.Errorf("environment variable %s is not set", envVar)
			}
			v.Field(i).SetString(value)
		}
	}

	return nil
}
