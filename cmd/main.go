package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Aditya-Bhalerao/cohereDB/db"
	"github.com/Aditya-Bhalerao/cohereDB/utils"
	"github.com/Aditya-Bhalerao/cohereDB/web"
)

func main() {
	configPath := flag.String("config", "", "Path to server config file (e.g., server0.toml)")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("Please provide a config file path using the -config flag")
	}

	utils.NewLogger()
	utils.Logger.Info().Msg("cohereDB server starting...")

	var config web.Config
	err := utils.LoadTomlConfig(&config, *configPath)
	if err != nil {
		utils.Logger.Fatal().Err(err).Msg("Failed to load config file")
		return
	}

	peerServers := make(map[int]string)
	for _, srv := range config.PeerServers {
		peerServers[srv.Shard] = srv.Addr
	}

	// Initialize the database
	dbInstance, err := db.NewDatabase(fmt.Sprintf("../data/db_%d", config.Server.Shard))
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to initialize database")
		return
	}

	// Create and start the server instance
	server := web.NewServer(dbInstance, config.Server.Addr, config.Server.Shard, config.Database.ShardCount, &peerServers)
	server.RegisterHandlers()

	utils.Logger.Info().Msgf("Starting server at %s...", config.Server.Addr)
	if err := server.Start(); err != nil {
		utils.Logger.Fatal().Err(err).Msg("Server failed")
	}
}
