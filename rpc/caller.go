package rpc

import (
	"errors"
	"sync"
	"time"

	"github.com/anteater2/bitmesh/message"
)

// Caller represents a caller service where remote functions are declared.
// It sends the call to callee over the network
// and captures the correpsponding return value.
type Caller struct {
	receiver *message.Receiver
	retChan  map[uint32]chan interface{}
	nextID   func() (uint32, error)
	freeID   func(uint32)
}

// NewCaller creates a new Caller
func NewCaller(port int) (*Caller, error) {
	var c Caller
	var err error
	c.nextID, c.freeID = makeIDGenerator()
	c.retChan = make(map[uint32]chan interface{})
	c.receiver, err = message.NewReceiver(port, func(v interface{}) {
		reply := v.(reply)
		if ret, prs := c.retChan[reply.ID]; prs {
			ret <- reply.Ret
		}
	})
	if err != nil {
		return nil, err
	}
	c.receiver.Register(reply{})
	message.Register(call{})
	return &c, nil
}

// RemoteFunc is the type returned by Declare
type RemoteFunc func(addr string, arg interface{}) (interface{}, error)

// Declare registers a return type and makes a function
// which sends a call to the specified address and block until return or timeout
// There must be a Caller at the specified address to process the call correctly.
func (c *Caller) Declare(arg interface{}, ret interface{}, timeout time.Duration) RemoteFunc {
	message.Register(arg)
	c.receiver.Register(ret)
	return func(addr string, arg interface{}) (interface{}, error) {
		id, err := c.nextID()
		if err != nil {
			return nil, err
		}
		defer c.freeID(id)

		// prepare a channel to receive return value
		ret := make(chan interface{}, 1)
		c.retChan[id] = ret
		defer delete(c.retChan, id)

		// send the call
		call := call{id, c.receiver.Addr(), arg}
		err = message.Send(addr, call)
		if err != nil {
			return nil, err
		}

		// wait for return or timeout
		select {
		case val := <-ret:
			return val, nil
		case <-time.After(timeout):
			return nil, errors.New("time out")
		}
	}
}

// Start starts the caller
func (c *Caller) Start() error {
	return c.receiver.Start()
}

// Stop stops the caller
func (c *Caller) Stop() {
	c.receiver.Stop()
	c.nextID, c.freeID = makeIDGenerator()
	c.retChan = make(map[uint32]chan interface{})
}

func makeIDGenerator() (func() (uint32, error), func(id uint32)) {
	var mutex sync.Mutex
	var counter uint32
	var usedID = make(map[uint32]bool)
	nextID := func() (uint32, error) {
		mutex.Lock()
		if len(usedID) == 1<<32 {
			mutex.Unlock()
			return 0, errors.New("out of call id")
		}
		for {
			if _, prs := usedID[counter]; !prs {
				break
			}
			counter++
		}
		newID := counter
		counter++
		mutex.Unlock()
		return newID, nil
	}
	freeID := func(id uint32) {
		mutex.Lock()
		delete(usedID, id)
		mutex.Unlock()
	}
	return nextID, freeID
}
