package chord

import (
	"errors"
	"flag"
	"log"
	"net"
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

// Init initializes the configs
func Init() error {
	config.addr = GetOutboundIP()
	flag.Uint64Var(
		&config.bits,
		"n",
		0,
		"Create a new chord ring with a keyspace of size 2^numBits",
	)

	flag.StringVar(
		&config.introducer,
		"c",
		"",
		"Create a new node and connect to the specified ring address",
	)

	flag.Parse()

	if config.bits == 0 {
		return errors.New("you must specify the keyspace size of the chord ring")
	}

	if config.bits > 63 {
		return errors.New("invalid keyspace; maximum keyspace size is 63") // Not really, but easier for now
	}
	config.isCreator = config.introducer == ""
	config.maxKey = 1 << config.bits
	config.numFingers = config.bits - 1
	config.callerPort = 2000
	config.calleePort = 2001
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
