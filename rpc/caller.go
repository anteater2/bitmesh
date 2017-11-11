package rpc

import (
	"errors"
	"time"

	"github.com/anteater2/bitmesh/message"
)

// Caller represents a caller service where remote functions are declared.
// It sends the call to callee over the network
// and captures the correpsponding return value.
type Caller struct {
	sender   *message.Sender
	receiver *message.Receiver
	retChan  map[int64]chan interface{}
}

// NewCaller creates a new Caller
func NewCaller(port int) (*Caller, error) {
	var c Caller
	var err error
	c.retChan = make(map[int64]chan interface{})
	c.sender = message.NewSender()
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
	c.sender.Register(call{})
	return &c, nil
}

// RemoteFunc is the type returned by Declare
type RemoteFunc func(addr string, arg interface{}) (interface{}, error)

// Declare registers a return type and makes a function
// which sends a call to the specified address and block until return or timeout
// There must be a Callee at the specified address to process the call correctly.
func (c *Caller) Declare(arg interface{}, ret interface{}, timeout time.Duration) RemoteFunc {
	c.sender.Register(arg)
	c.receiver.Register(ret)
	return func(addr string, arg interface{}) (interface{}, error) {
		id := time.Now().Unix()

		// prepare a channel to receive return value
		ret := make(chan interface{}, 1)
		c.retChan[id] = ret
		defer delete(c.retChan, id)

		// send the call
		call := call{id, c.receiver.Addr(), arg}
		err := c.sender.Send(addr, call)
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

// Addr returns the address of Caller (only valid when Caller is running)
func (c *Caller) Addr() string {
	return c.receiver.Addr()
}

// Stop stops the caller
func (c *Caller) Stop() {
	c.receiver.Stop()
	c.retChan = make(map[int64]chan interface{})
}
