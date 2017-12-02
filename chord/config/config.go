package config

import (
	"errors"
	"flag"
	"log"
	"net"
)

var (
	addr       string
	bits       uint64
	callerPort uint16 = 2000
	calleePort uint16 = 2001
	introducer string
	isCreator  bool
	maxKey     uint64
	numFingers uint64
	username   string
)

// GetOutboundIP gets preferred outbound IP of this machine using a filthy hack
// The connection should not actually require the Google DNS service (the 8.8.8.8),
// but by creating it we can see what our preferred IP is.
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// CallerPort returns the port of the caller
func CallerPort() uint16 {
	return callerPort
}

// CalleePort returns the port of the callee
func CalleePort() uint16 {
	return calleePort
}

// Creator returns the truth value of whether this node is the first node
// in a chord network
func Creator() bool {
	return isCreator
}

// Addr is this node's preferred outbound IP
func Addr() string {
	return addr
}

// Init initializes the configs
func Init() error {
	addr = GetOutboundIP()
	flag.Uint64Var(
		&bits,
		"n",
		0,
		"Create a new chord ring with a keyspace of size 2^numBits",
	)

	flag.StringVar(
		&introducer,
		"c",
		"",
		"Create a new node and connect to the specified ring address",
	)

	flag.StringVar(
		&username,
		"u",
		"",
		"the username to use",
	)

	flag.Parse()

	if bits == 0 {
		return errors.New("you must specify the keyspace size of the chord ring")
	}

	if bits > 63 {
		return errors.New("invalid keyspace; maximum keyspace size is 63") // Not really, but easier for now
	}
	isCreator = introducer == ""
	maxKey = 1 << bits
	numFingers = bits - 1
	return nil
}

// Introducer returns the introducing address
func Introducer() string {
	return introducer
}

// MaxKey returns the size of the key space.
func MaxKey() uint64 {
	return maxKey
}

// NumFingers returns the size of a finger table
func NumFingers() uint64 {
	return numFingers
}

// Username returns the username of the node
func Username() string {
	return username
}
