package dht

import (
	"fmt"
	"log"
	"time"

	"../../chord-node/config"
	"../../chord-node/start"
	"github.com/anteater2/bitmesh/rpc"
	"github.com/anteater2/chord-node/key"
)

type request struct {
	key key.Key
}

// Init the distributed hash table
func Init() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Creating local node @IP%s on its own ring of size %d...\n", config.Addr(), config.MaxKey())
	start.CreateLocalNode()
	go start.RPCCallee.Start()
	go start.RPCCaller.Start()
	callee, _ := rpc.NewCallee(2003)
	callee.Implement(receiveResponse)

	// if !config.Creator() {
	// 	start.Join(config.Introducer())
	// }
	fmt.Printf("Beginning stabilizer...\n")
	go start.Stabilize()
	go start.FixFingers()
	select {}
}

//Put a key into the distributed hash table,
// assume the remote node the key represent is a initilized hash table
//joins the two hashtable together
func Put(key key.Key, address string) {
	// remote := node.RemoteNode{address, key}

	start.Join(address)

}

//Get the addreess of a remote node by the key,
// return the address of the remote node,
//if the key is not in the hash table or
//is no longer alive then return null;
func Get(key key.Key) string {
	address := "local" // should be the address of in the finger table
	caller, _ := rpc.NewCaller(2001)
	// callee, _ := rpc.NewCallee(2002)
	sendReq := caller.Declare(request{}, 0, time.Second)
	res, err := sendReq(address, request{key}) //probably need to consider the port number but not sure
	// callee.Implement(receiveResponse)
	go caller.Start()
	// go callee.Start()

}

//this was meant to be implemented by the rpc callee
//receive a key then if the node is in charge of the
//the key, return the address; else send the information to
// to node that is responisble
func receiveRequest(key string) {

}

func receiveResponse(address string) {

}
