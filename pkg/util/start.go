// Package util implements different utilities required by the EdgeCluster service
package util

import (
	"log"
	"os"
	"os/signal"

	"github.com/decentralized-cloud/edge-core/services/configuration"
	"github.com/decentralized-cloud/edge-core/services/cron/ipgeolocation"
	"github.com/decentralized-cloud/edge-core/services/transport/http"
	"go.uber.org/zap"
)

var configurationService configuration.ConfigurationContract

// StartService setups all dependecies required to start the EdgeCluster service and
// start the service
func StartService() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	if err = setupDependencies(logger); err != nil {
		logger.Fatal("Failed to setup dependecies", zap.Error(err))
	}

	geolocationUpdaterService, err := ipgeolocation.NewCronService(
		logger,
		configurationService)
	if err != nil {
		logger.Fatal("Failed to create Geolocation Updater service", zap.Error(err))
	}

	httpTansportService, err := http.NewTransportService(
		logger,
		configurationService)
	if err != nil {
		logger.Fatal("Failed to create HTTP transport service", zap.Error(err))
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)

	if configurationService.ShouldUpdatePublciIPAndGeolocationDetails() {
		go func() {
			if serviceErr := geolocationUpdaterService.Start(); serviceErr != nil {
				logger.Fatal("Failed to start Geolocation Updater service", zap.Error(serviceErr))
			}
		}()
	}

	go func() {
		if serviceErr := httpTansportService.Start(); serviceErr != nil {
			logger.Fatal("Failed to start HTTP transport service", zap.Error(serviceErr))
		}
	}()

	go func() {
		<-signalChan
		logger.Info("Received an interrupt, stopping services...")

		if configurationService.ShouldUpdatePublciIPAndGeolocationDetails() {
			if err := geolocationUpdaterService.Stop(); err != nil {
				logger.Error("Failed to stop Geolocation Updater service", zap.Error(err))
			}
		}

		if err := httpTansportService.Stop(); err != nil {
			logger.Error("Failed to stop HTTP transport service", zap.Error(err))
		}

		close(cleanupDone)
	}()
	<-cleanupDone
}

func setupDependencies(logger *zap.Logger) (err error) {
	if configurationService, err = configuration.NewEnvConfigurationService(); err != nil {
		return
	}

	return
}
