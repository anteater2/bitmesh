package dht

import (
	"github.com/anteater2/bitmesh/chord"
)

// DHT represents a client for a distributed hash table.
type DHT struct {
	node   string
	bits   uint64
	caller *chord.NodeCaller
}

// New creates a client to access DHT
func New(node string, receivePort uint16, bits uint64) (*DHT, error) { // configuration
	caller, err := chord.NewNodeCaller(receivePort)
	if err != nil {
		return nil, err
	}
	return &DHT{
		node:   node,
		bits:   bits,
		caller: caller,
	}, nil
}

// Start ...
func (dht *DHT) Start() {
	dht.caller.Start()
}

// Put puts a key-value pair into dht.
func (dht *DHT) Put(k string, v string) error {
	hashk := chord.Hash(k, 1<<dht.bits)
	remote, err := dht.caller.FindSuccessor(dht.node, hashk)
	if err != nil {
		return err
	}
	err = dht.caller.Put(remote.Address, k, []byte(v))
	if err != nil {
		return err
	}
	return nil
}

// Get gets the value corresponding to the key from dht
func (dht *DHT) Get(k string) (string, error) {
	hashk := chord.Hash(k, 1<<dht.bits)
	remote, err := dht.caller.FindSuccessor(dht.node, hashk)
	if err != nil {
		return "", err
	}
	v, err := dht.caller.Get(remote.Address, k)
	if err != nil {
		return "", err
	}
	return string(v), nil
}
