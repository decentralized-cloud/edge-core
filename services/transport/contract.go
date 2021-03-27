// Package transport implements different transport services required by the project service
package transport

// TransportContract declares the methods to be implemented by the transport service
type TransportContract interface {
	// Start the transport service.
	// Returns error if something goes wrong.
	Start() error

	// Stop the transport service.
	// Returns error if something goes wrong.
	Stop() error
}
