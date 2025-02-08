package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arbha1erao/cohereDB/cohere/db"
	"github.com/arbha1erao/cohereDB/cohere/utils"
	"github.com/arbha1erao/cohereDB/cohere/web"
)

func main() {
	configPath := flag.String("config", "", "Path to server config file (e.g., server0.toml)")
	cleanup := flag.Bool("cleanup", false, "Clean up the database when server shuts down")
	container := flag.Bool("container", false, "CohereDB running inside a container")
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
	dbPath := fmt.Sprintf("../data/db_%d", config.Server.Shard)
	dbInstance, err := db.NewDatabase(dbPath)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to initialize database")
		return
	}

	var serverAddr string
	if *container {
		utils.Logger.Info().Msgf("Running in container mode. Using container address: %s", config.Server.ContainerAddr)
		serverAddr = config.Server.ContainerAddr
	} else {
		utils.Logger.Info().Msgf("Running in host mode. Using address: %s", config.Server.Addr)
		serverAddr = config.Server.Addr
	}

	// Create and start the server instance
	server := web.NewServer(dbInstance, serverAddr, config)
	server.RegisterHandlers()

	// Setup graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		utils.Logger.Info().Msgf("Starting server at %s...", serverAddr)
		if err := server.Start(); err != nil {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		utils.Logger.Fatal().Err(err).Msg("Server failed")
	case sig := <-signalChan:
		utils.Logger.Info().Msgf("Received signal %s, initiating shutdown...", sig)

		if *cleanup {
			utils.Logger.Info().Msg("Cleanup flag set, cleaning up database...")
			if err := dbInstance.Cleanup(); err != nil {
				utils.Logger.Error().Err(err).Msg("Failed to cleanup database")
			} else {
				utils.Logger.Info().Msg("Database cleanup completed")
			}
		} else {
			if err := dbInstance.Close(); err != nil {
				utils.Logger.Error().Err(err).Msg("Error closing database")
			}
		}

		utils.Logger.Info().Msg("Server shutdown completed")
	}
}
