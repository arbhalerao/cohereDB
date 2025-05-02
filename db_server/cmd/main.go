package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/arbha1erao/cohereDB/db"
	"github.com/arbha1erao/cohereDB/db_server"
	grpc_server "github.com/arbha1erao/cohereDB/db_server/grpc"
	http_server "github.com/arbha1erao/cohereDB/db_server/http"
	"github.com/arbha1erao/cohereDB/utils"
)

type Config struct {
	Server struct {
		Region       string `toml:"region"`
		HTTP_Addr    string `toml:"http_addr"`
		GRPC_Addr    string `toml:"grpc_addr"`
		MANAGER_Addr string `toml:"manager_addr"`
	} `toml:"server"`
}

func main() {
	configPath := flag.String("config", "config.toml", "Path to the config file")
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
	httpAddr := config.Server.HTTP_Addr
	grpcAddr := config.Server.GRPC_Addr
	managerAddr := config.Server.MANAGER_Addr

	// Initialize database
	dbPath := fmt.Sprintf("../../data/db_%s", region)
	database, err := db.NewDatabase(dbPath)
	if err != nil {
		utils.Logger.Fatal().Err(err).Msg("Failed to initialize database")
		return
	}
	defer database.Close()

	// Initialize HTTP Server
	httpServer := http_server.NewServer(database, httpAddr)
	httpServer.RegisterHandlers()

	// Initialize gRPC Server
	grpcService := grpc_server.NewServer(database, grpcAddr)

	// Use DBManagerClient for registration
	ready := make(chan bool)
	managerClient := db_server.NewDBManagerClient(managerAddr, region)
	go managerClient.RegisterWithManager(region, httpAddr, grpcAddr, ready)

	// Wait for registration to complete before proceeding
	utils.Logger.Info().Msg("Waiting for registration with db_manager...")
	<-ready
	close(ready)
	utils.Logger.Info().Msg("Registration successful. Starting servers...")

	// Handle shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)

	// Start HTTP server
	go func() {
		defer wg.Done()
		utils.Logger.Info().Msgf("Starting HTTP server on %s", httpAddr)
		if err := httpServer.Start(); err != nil {
			utils.Logger.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

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

	if err := httpServer.Shutdown(); err != nil {
		utils.Logger.Fatal().Err(err).Msg("HTTP server shutdown error")
	}

	wg.Wait()

	utils.Logger.Info().Msg("Servers stopped successfully.")
}
