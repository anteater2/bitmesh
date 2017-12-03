package dht

import (
	"github.com/anteater2/bitmesh/chord"
)

// Start starts dht service
func Start() {
	chord.Start()
}

// Put puts a key-value pair into dht.
func Put(key string, value string) error {
	return chord.Put(key, []byte(value))
}

// Get gets the value corresponding to the key from dht
func Get(key string) (string, error) {
	bytes, err := chord.Get(key)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
