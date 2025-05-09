package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/arbha1erao/cohereDB/db"
	"github.com/arbha1erao/cohereDB/db_server/db_manager_client"
	grpc_server "github.com/arbha1erao/cohereDB/db_server/grpc"
	"github.com/arbha1erao/cohereDB/utils"
)

type Config struct {
	Server struct {
		Region       string `toml:"region"`
		GRPC_Addr    string `toml:"grpc_addr"`
		MANAGER_Addr string `toml:"manager_addr"`
	} `toml:"server"`
}

func main() {
	configPath := flag.String("config", "config.toml", "Path to the config file")
	register := flag.Bool("register", false, "Indicates if registration should happen")
	flag.Parse()

	utils.NewLogger()
	utils.Logger.Info().Msg("cohereDB server starting...")

	// Load configuration
	var config Config
	err := utils.LoadTomlConfig(&config, *configPath)
	if err != nil {
		utils.Logger.Fatal().Err(err).Msg("Failed to load config file")
		return
	}

	region := config.Server.Region
	grpcAddr := config.Server.GRPC_Addr
	managerAddr := config.Server.MANAGER_Addr

	// Initialize database
	dbPath := fmt.Sprintf("../../data/db_%s", region)
	database, err := db.NewDatabase(dbPath)
	if err != nil {
		utils.Logger.Fatal().Err(err).Msg("Failed to initialize database")
		return
	}
	defer func() {
		if err := database.Close(); err != nil {
			utils.Logger.Error().Err(err).Msg("Failed to close database")
		}
	}()

	// Initialize gRPC Server
	grpcService := grpc_server.NewServer(database, grpcAddr)

	// Use DBManagerClient for registration, if registration flag is set
	ready := make(chan bool)
	if *register {
		managerClient := db_manager_client.NewDBManagerClient(managerAddr, region)
		go managerClient.RegisterWithManager(region, grpcAddr, ready)

		// Wait for registration to complete before proceeding
		utils.Logger.Info().Msg("Waiting for registration with db_manager...")
		<-ready
		utils.Logger.Info().Msg("Registration successful. Starting servers...")
	}
	close(ready)

	// Handle shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	// Start gRPC server
	go func() {
		defer wg.Done()
		utils.Logger.Info().Msgf("Starting gRPC server on %s", grpcAddr)
		if err := grpcService.Start(); err != nil {
			utils.Logger.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	<-stop
	utils.Logger.Info().Msg("Shutting down servers...")

	grpcService.Stop()

	wg.Wait()

	utils.Logger.Info().Msg("Servers stopped successfully.")
}
