package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rabbitprincess/x402-facilitator/api"
	"github.com/rabbitprincess/x402-facilitator/facilitator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "facilitator",
	Short: "Start the facilitator server",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var (
	configPath string
)

func init() {
	serverCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.toml", "Path to the configuration file")
}

func main() {
	if err := serverCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute command")
	}
}

func run() {
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration, Shutting down...")
	}
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	facilitator, err := facilitator.NewFacilitator(&logger, config.Scheme, config.Url, config.PrivateKey)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create facilitator, Shutting down...")
	}

	// &logger, facilitator 넣어서 내부에서 써주세영
	api := api.NewServer()

	// Initialize Server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: api.Router(),
	}

	go func() {
		logger.Info().Msgf("Starting server on port %d", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start server, Shutting down...")
		}
	}()

	// Graceful shutdown handling
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Failed to shutdown server gracefully")
	}
	logger.Info().Msg("Server shutdown gracefully")
}
