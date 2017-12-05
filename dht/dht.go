package dht

import (
	"fmt"
	"time"

	"github.com/anteater2/bitmesh/chord"
	"github.com/anteater2/bitmesh/chord/key"
	"github.com/anteater2/bitmesh/rpc"
)

// DHT ...
type DHT struct {
	node      string
	chordPort uint16
	caller    *rpc.Caller
	rpcGet    rpc.RemoteFunc
	rpcPut    rpc.RemoteFunc
	rpcLookup rpc.RemoteFunc
}

// New creates a client to access DHT
func New(node string, chordPort uint16, receivePort uint16) (*DHT, error) { // configuration
	caller, err := rpc.NewCaller(receivePort)
	if err != nil {
		return nil, err
	}
	return &DHT{
		node:      node,
		chordPort: chordPort,
		caller:    caller,
		rpcGet:    caller.Declare("", chord.GetKeyResponse{}, 3*time.Second),
		rpcPut:    caller.Declare(chord.PutKeyRequest{}, true, 3*time.Second),
		rpcLookup: caller.Declare(key.Key(0), chord.RemoteNode{}, 5*time.Second),
	}, nil
}

// Start ...
func (dht *DHT) Start() {
	dht.caller.Start()
}

// Put puts a key-value pair into dht.
func (dht *DHT) Put(k string, v string) error {
	hashk := key.Hash(k, 1<<10)
	request := chord.PutKeyRequest{KeyString: k, Data: []byte(v)}
	remoteNode, err := dht.rpcLookup(dht.node, hashk)
	if err != nil {
		return err
	}
	remote := remoteNode.(chord.RemoteNode).Address
	ok, err := dht.rpcPut(joinAddrPort(remote, dht.chordPort), request)
	if err != nil {
		return err
	}
	if !ok.(bool) {
		return fmt.Errorf("put failed")
	}
	return nil
}

// Get gets the value corresponding to the key from dht
func (dht *DHT) Get(k string) (string, error) {
	hashk := key.Hash(k, 1<<10)
	remoteNode, err := dht.rpcLookup(dht.node, hashk)
	if err != nil {
		return "", err
	}
	remote := remoteNode.(chord.RemoteNode).Address
	res, err := dht.rpcGet(joinAddrPort(remote, dht.chordPort), k)
	response := res.(chord.GetKeyResponse)
	if err != nil {
		return "", err
	}
	if response.Error == false {
		return "", fmt.Errorf("get failed")
	}
	return string(response.Data), nil
}

func joinAddrPort(addr string, port uint16) string {
	return fmt.Sprintf("%s:%d", addr, port)
}
