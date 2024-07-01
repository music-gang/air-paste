package airpaste

import "errors"

// RandomString returns a securely generated random string.
var RandomString func(n int) (string, error)

// Gateway directs the requests to the approriate handlers and provides the necessary data to them.
type Gateway struct {
	kv AirDatastore
}

// NewGateway creates a new Gateway.
func NewGateway(kv AirDatastore) *Gateway {
	return &Gateway{
		kv: kv,
	}
}

// GetHandler returns the value for the given key. The second return value is true if the key exists, and false otherwise.
func (g *Gateway) GetHandler(key string) (string, bool) {
	return g.kv.Get(key)
}

// SetHandler sets the value and returns the generated key.
func (g *Gateway) SetHandler(value string, opt ...SetOptions) (key string, err error) {

	size := 4

	for {
		if size > 8 {
			return "", errors.New("could not generate a unique key")
		}

		key, err = RandomString(size)
		if err != nil {
			return "", err
		}

		if _, ok := g.kv.Get(key); !ok {
			g.kv.Set(key, value, opt...)
			break
		}

		size++
	}

	return
}
