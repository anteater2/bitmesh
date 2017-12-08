package main

import (
	"flag"
	"log"
	"net"

	"github.com/anteater2/bitmesh/chord"
)

func main() {
	var bits uint64
	var introducer string
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

	flag.Parse()
	err := chord.Start(getOutboundIP(), 2001, 2000, bits)
	if err != nil {
		panic(err)
	}
	if introducer != "" {
		err = chord.Join(introducer)
		if err != nil {
			panic(err)
		}
	}
	select {}
}

// getOutboundIP gets preferred outbound IP of this machine using a filthy hack
// The connection should not actually require the Google DNS service (the 8.8.8.8),
// but by creating it we can see what our preferred IP is.
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
