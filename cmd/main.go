package main

import (
	"github.com/Aditya-Bhalerao/cohereDB/db"
	utils "github.com/Aditya-Bhalerao/cohereDB/utils"
)

type ServiceConfig struct {
	HTTP_ADDR string `env:"HTTP_ADDR"`
}

type BadgerConfig struct {
	BADGER_DB_PATH string `env:"BADGER_DB_PATH"`
}

func main() {
	utils.NewLogger()

	utils.Logger.Info().Msg("cohereDB starting...")

	var serviceConfig ServiceConfig
	var badgerConfig BadgerConfig

	err := utils.LoadConfig(&serviceConfig, &badgerConfig)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to load env config")
		return
	}

	_, err = db.NewDatabase(badgerConfig.BADGER_DB_PATH)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to initialize database")
		return
	}
}
