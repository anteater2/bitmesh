package dht

import (
	"fmt"
	"log"
	"time"

	"../../chord-node/config"
	"../../chord-node/key"
	"../../chord-node/start"
	"github.com/anteater2/bitmesh/rpc"
)

//use port 2001 to send request
//use port 2002 to receive response
//use port 2003 to receive request
//use port 2004 to send response
//use port 2005 to act as the intermediate needed?

//Request is sent when call get
//key: the key of the node that is looking for
//Requester : the ip address of the requester so the node in charge
//of key can directly send to data back to the requester
type Request struct {
	key       key.Key
	Requester string
}

//Response is sent when the corresponding node return the address
type Response struct {
	address string
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
	callee.Implement(receiveRequest)
	go callee.Start()
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
	remote := start.FindSuccessor(key)
	address := remote.Address // should be the address of the successor of the node
	//ToDo check the address is its own address or not
	caller, _ := rpc.NewCaller(2001)
	sendReq := caller.Declare(Request{}, 0, time.Second)
	res, err := sendReq(address+"2003", Request{key, "127.0.0.1"}) //use 127.0.0.1 right now to act as placeholder
	if err != nil {
		panic(err)
	}
	go caller.Start()
	fmt.Printf("the address is %s\n", res.(string))
	return res.(string)
}

//this was meant to be implemented by the  callee
//the callee is supposed to be the successor of the
//node we are looking for
//then return the address of the node to the caller
func receiveRequest(request Request) string {
	remote := start.GetPredecessor(1)
	address := remote.Address

	return address
}

// func receiveResponse(address Response) {

// }
