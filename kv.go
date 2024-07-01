package airpaste

import "time"

// SetOptions contains the options for setting a key-value pair.
type SetOptions struct {
	// TTL is the time-to-live for the key-value pair
	TTL *time.Duration
}

// AirDatastore is the interface that wraps the basic Get and Set methods.
type AirDatastore interface {
	// Get returns the value for the given key. The second return value is true if the key exists, and false otherwise.
	Get(key string) (string, bool)
	// Set sets the value for the given key, with the given time-to-live.
	Set(key, value string, opt ...SetOptions)
}
