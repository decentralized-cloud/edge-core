// Package configuration implements configuration service required by the edge-core service
package configuration

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	commonErrors "github.com/micro-business/go-core/system/errors"
)

type envConfigurationService struct {
}

// NewEnvConfigurationService creates new instance of the EnvConfigurationService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewEnvConfigurationService() (ConfigurationContract, error) {
	return &envConfigurationService{}, nil
}

// GetHttpHost returns HTTP host name
// Returns the HTTP host name
func (service *envConfigurationService) GetHttpHost() string {
	return os.Getenv("HTTP_HOST")
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

// GetEdgeClusterType returns the type of edge cluster such as K3S
// Returns the type of edge cluster or error if something goes wrong
func (service *envConfigurationService) GetEdgeClusterType() (ClusterType, error) {
	switch value := strings.Trim(os.Getenv("EDGE_CLUSTER_TYPE"), " "); value {
	case "K3S":
		return K3S, nil
	case "":
		return Unknown, commonErrors.NewUnknownError(
			"EDGE_CLUSTER_TYPE is required")
	default:
		return Unknown, commonErrors.NewUnknownError(
			fmt.Sprintf("Could not figure out the edge cluster type from the given EDGE_CLUSTER_TYPE (%s)", value))
	}
}

// ShouldUpdatePublciIPAndGeolocationDetails determines whether the edge-core should periodically check
// for the node public IP address and geolocation details
// Returns true if the edge-core should periodically check for the node public IP address and
// geolocation details otherwise returns false
func (service *envConfigurationService) ShouldUpdatePublciIPAndGeolocationDetails() bool {
	if value := strings.Trim(os.Getenv("UPDATE_PUBLIC_IP_GEOLOCATION_DETAILS"), " "); value == "true" {
		return true
	}

	return false
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
