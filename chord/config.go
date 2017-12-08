package chord

import (
	"errors"
)

var config struct {
	addr       string
	bits       uint64
	callerPort uint16
	calleePort uint16
	introducer string
	isCreator  bool
	maxKey     uint64
	numFingers uint64
}

// Init initializes the configs
func Init(addr string, calleePort uint16, callerPort uint16, bits uint64) error {
	config.addr = addr
	config.calleePort = calleePort
	config.callerPort = callerPort
	config.bits = bits

	if config.bits > 63 {
		return errors.New("invalid keyspace; maximum keyspace size > 63")
	}
	config.isCreator = config.introducer == ""
	config.maxKey = 1 << config.bits
	config.numFingers = config.bits - 1
	return nil
}

// Introducer returns the introducing address
func Introducer() string {
	return config.introducer
}

// MaxKey returns the size of the key space.
func MaxKey() uint64 {
	return config.maxKey
}

// NumFingers returns the size of a finger table
func NumFingers() uint64 {
	return config.numFingers
}
