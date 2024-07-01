package kv

import (
	"sync"
	"time"

	airpaste "github.com/music-gang/air-paste"
)

type Key = string

// defaultTTL is the default time-to-live for a key-value pair
const defaultTTL = 2 * time.Minute

// Data is the type for a key-value pair in the datastore.
// It contains the actual value, the time-to-live, and the creation time.
// It is not synchronized, so it should be used only within a synchronized context.
type Data struct {
	// value is the actual value of the key-value pair
	value string
	// ttl is the time-to-live for the key-value pair
	// A negative value means the key-value pair never expires
	ttl time.Duration
	// crt is the creation time for the key-value pair
	// is allowed to be zero if ttl is negative, otherwise it is considered as actual creation time
	crt time.Time
}

// NewData creates a new Data with the given value, time-to-live, and creation time.
func NewData(value string, ttl time.Duration, crt time.Time) Data {
	return Data{
		value: value,
		ttl:   ttl,
		crt:   crt,
	}
}

// KV is a simple in-memory key-value store.
type KV struct {
	// data contains the actual key-value pairs
	data map[Key]Data
}

// NewDatastore creates a new Datastore.
func NewDatastore() *KV {
	return &KV{
		data: make(map[Key]Data),
	}
}

// Get returns the value for the given key. The second return value is true if the key exists, and false otherwise.
// Also, it checks if the key has expired, and if so, it deletes it and returns false.
func (d *KV) Get(key string) (string, bool) {

	// check if the key passed existss in the datastore
	data, ok := d.data[key]
	if !ok {
		return "", false
	}

	if data.ttl >= 0 {
		// check if the key has expired
		if time.Since(data.crt) > data.ttl {
			del(d, key)
			return "", false
		}
	}

	return data.value, true
}

// Set sets the value for the given key, with the given time-to-live.
func (d *KV) Set(key, value string, opt ...airpaste.SetOptions) {

	var options airpaste.SetOptions

	if len(opt) > 0 {
		options = opt[0]
	}

	var ttl time.Duration
	if options.TTL != nil {
		ttl = *options.TTL
	} else {
		ttl = defaultTTL
	}

	d.data[key] = NewData(value, ttl, time.Now().UTC())
}

// del deletes the key-value pair with the given key from the datastore.
// It is a helper function to delete a key-value pair from the datastore and should be unexported outside of this package.
func del(d *KV, key string) {
	delete(d.data, key)
}

// SyncedKV is a synchronized version of Datastore.
// It embeds a Datastore and adds a mutex to synchronize access to the underlying Datastore.
type SyncedKV struct {
	*KV
	mux sync.RWMutex
}

// NewSyncedDatastore creates a new SyncedDatastore.
func NewSyncedDatastore() *SyncedKV {
	return &SyncedKV{
		KV: NewDatastore(),
	}
}

// Get returns the value for the given key. The second return value is true if the key exists, and false otherwise.
// Also, it checks if the key has expired, and if so, it deletes it and returns false.
func (d *SyncedKV) Get(key string) (string, bool) {
	d.mux.RLock()
	defer d.mux.RUnlock()
	return d.KV.Get(key)
}

// Set sets the value for the given key, with the given time-to-live.
func (d *SyncedKV) Set(key, value string, opt ...airpaste.SetOptions) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.KV.Set(key, value, opt...)
}
