// Package configuration implements configuration service required by the edge-core service
package configuration

import (
	commonErrors "github.com/micro-business/go-core/system/errors"
	"os"
	"strconv"
	"strings"
)

type envConfigurationService struct {
}

// NewEnvConfigurationService creates new instance of the EnvConfigurationService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewEnvConfigurationService() (ConfigurationContract, error) {
	return &envConfigurationService{}, nil
}

// GetHttpHost returns HTTP host name
// Returns the HTTP host name or error if something goes wrong
func (service *envConfigurationService) GetHttpHost() (string, error) {
	return os.Getenv("HTTP_HOST"), nil
}

// GetHttpPort returns HTTP port number
// Returns the HTTP port number or error if something goes wrong
func (service *envConfigurationService) GetHttpPort() (int, error) {
	valueStr := os.Getenv("HTTP_PORT")
	if strings.Trim(valueStr, " ") == "" {
		return 0, commonErrors.NewUnknownError("HTTP_PORT is required")
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, commonErrors.NewUnknownErrorWithError("Failed to convert HTTP_PORT to integer", err)
	}

	return value, nil
}

// GetRunningNodeName returns the name of the node that currently running the pod
// Returns the name of the node that currently running the pod or error if something goes wrong
func (service *envConfigurationService) GetRunningNodeName() (string, error) {
	value := os.Getenv("NODE_NAME")
	if strings.Trim(value, " ") == "" {
		return "", commonErrors.NewUnknownError("NODE_NAME is required")
	}

	return value, nil
}

// GetGeolocationUpdaterCronSpec returns Geolocation Updater updating interval
// Returns the Geolocation Updater updating interval or error if something goes wrong
func (service *envConfigurationService) GetGeolocationUpdaterCronSpec() (string, error) {
	value := os.Getenv("GEOLOCATION_UPDATER_CRON_SPEC")
	if strings.Trim(value, " ") == "" {
		return "", commonErrors.NewUnknownError("GEOLOCATION_UPDATER_CRON_SPEC is required")
	}

	return value, nil
}

// GetIpinfoUrl returns the URL to the Ipinfo website that returns the node public IP address
// Returns the URL to the Ipinfo website that returns the node public IP address or error if something goes wrong
func (service *envConfigurationService) GetIpinfoUrl() (string, error) {
	value := os.Getenv("IPINFO_URL")
	if strings.Trim(value, " ") == "" {
		return "https://ipinfo.io", nil
	}

	return value, nil
}

// GetIpinfoAccessToken returns the access token to be used when making request to the Ipinfo website
// to return the node public IP address
// Returns the access token to be used when making request to the Ipinfo website to return the node
// public IP address or error if something goes wrong
func (service *envConfigurationService) GetIpinfoAccessToken() (string, error) {
	return os.Getenv("IPINFO_ACCESS_TOKEN"), nil
}
