// Package cron implements different cron services required by the edge-core
package cron

// CronContract declares the methods to be implemented by the cron service
type CronContract interface {
	// Start the cron service.
	// Returns error if something goes wrong.
	Start() error

	// Stop the cron service.
	// Returns error if something goes wrong.
	Stop() error
}
