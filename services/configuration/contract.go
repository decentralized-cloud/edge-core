// Package configuration implements configuration service required by the edge-core service
package configuration

// ConfigurationContract declares the service that provides configuration required by different Tenat modules
type ConfigurationContract interface {
	// GetHttpHost returns HTTP host name
	// Returns the HTTP host name or error if something goes wrong
	GetHttpHost() (string, error)

	// GetHttpPort returns HTTP port number
	// Returns the HTTP port number or error if something goes wrong
	GetHttpPort() (int, error)

	// GetRunningNodeName returns the name of the node that currently running the pod
	// Returns the name of the node that currently running the pod or error if something goes wrong
	GetRunningNodeName() (string, error)

	// GetGeolocationUpdaterCronSpec returns Geolocation Updater updating interval
	// Returns the Geolocation Updater updating interval or error if something goes wrong
	GetGeolocationUpdaterCronSpec() (string, error)

	// GetIpinfoUrl returns the URL to the Ipinfo website that returns the node public IP address
	// Returns the URL to the Ipinfo website that returns the node public IP address or error if something goes wrong
	GetIpinfoUrl() (string, error)

	// GetIpinfoAccessToken returns the access token to be used when making request to the Ipinfo website
	// to return the node public IP address
	// Returns the access token to be used when making request to the Ipinfo website to return the node
	// public IP address or error if something goes wrong
	GetIpinfoAccessToken() (string, error)
}
