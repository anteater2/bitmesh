package chord

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/anteater2/bitmesh/chord/config"
	"github.com/anteater2/bitmesh/rpc"
)

var internalTable *HashTable

var my struct {
	key Key
}

var Address string
var fingers []*RemoteNode
var predecessor *RemoteNode
var successor *RemoteNode
var rpcCaller *rpc.Caller
var rpcCallee *rpc.Callee

var rpcFindSuccessor rpc.RemoteFunc
var rpcNotify rpc.RemoteFunc
var rpcGetPredecessor rpc.RemoteFunc
var rpcIsAlive rpc.RemoteFunc
var rpcPutKey rpc.RemoteFunc
var rpcGetKey rpc.RemoteFunc
var rpcGetKeyRange rpc.RemoteFunc

// rpcPutKeyBackup is used to backup a key to the node's predecessor.  This way, if the node fails, the key is duplicated.
var rpcPutKeyBackup rpc.RemoteFunc

// RemoteNode holds information for connecting to a remote node
type RemoteNode struct {
	Address string
	Key     Key
}

type GetKeyResponse struct {
	Data  []byte
	Error bool
}

type PutKeyRequest struct {
	KeyString string
	Data      []byte
}

type GetKeyRangeRequest struct {
	Start Key
	End   Key
}

// ClosestPrecedingNode finds the closest preceding node to the key in this node's finger table.
// This doesn't need any RPC.
func ClosestPrecedingNode(key Key) RemoteNode {
	for i := config.NumFingers() - 1; i > 0; i-- { // WARNING: GO DOES THIS i>0 CHECK AT THE END OF THE LOOP!
		//log.Printf("Checking finger %d\n", i)
		if fingers[i] == nil {
			panic("You attempted to find ClosestPrecedingNode without an initialized finger table!")
		}
		if fingers[i].Key.BetweenExclusive(my.key, key) {
			return *fingers[i]
		}
	}
	return RemoteNode{Address: Address, Key: my.key}
}

// FindSuccessor finds the successor node to the key.  This may require RPC calls.
func FindSuccessor(key Key) RemoteNode {
	if key.BetweenEndInclusive(my.key, successor.Key) {
		// key is between this node and its successor
		return *successor
	}
	target := ClosestPrecedingNode(key)
	if target.Address == Address {
		log.Printf("[DIAGNOSTIC] Infinite loop detected!\n")
		log.Printf("[DIAGNOSTIC] This is likely because of a bad finger table.\n")
		panic("This is probably a serious bug.")
	}
	// Now, we have to do an RPC on target to find the successor.

	interf, err := rpcFindSuccessor(joinAddrPort(target.Address, config.CalleePort()), key)
	if err != nil {
		log.Printf("[DIAGNOSTIC] Remote target is " + joinAddrPort(target.Address, config.CalleePort()) + "\n")
		log.Print(err)
		panic("rpcFindSuccessor failed!")
	}
	rv := interf.(RemoteNode)
	return rv
}

// Notify notifies the successor that you are the predecessor
func Notify(node RemoteNode) int {
	if predecessor == nil || node.Key.BetweenExclusive(predecessor.Key, my.key) {
		log.Printf("Got notify from %s!  New predecessor: %d\n", node.Address, node.Key)
		predecessor = &node
		if predecessor.Address != Address {

			rvInterf, err := rpcGetKeyRange(joinAddrPort(predecessor.Address, config.CalleePort()), GetKeyRangeRequest{my.key, predecessor.Key})
			if err != nil {
				log.Fatal(err)
			}
			rv := rvInterf.([]HashEntry)
			for _, entry := range rv {
				internalTable.Put(entry.Key, entry.Value)
			}
		}
	}
	return 0 // Necessary to interface with rpcCaller
}

//Stabilize the Successor and Predecessor fields of this node.
//This is a goroutine and never terminates.
func Stabilize() {
	for true { // This is how while loops work.  Not even joking.
		var remote RemoteNode
		if predecessor == nil {
			log.Printf("Null predecessor!  New predecessor: %d\n", successor.Key)
			predecessor = successor
		}
		if successor.Address == Address {
			// Avoid making an RPC call to ourselves
			remote = *predecessor
		} else {

			remoteInterf, err := rpcGetPredecessor(joinAddrPort(successor.Address, config.CalleePort()), 0) // 0 is a dummy value so that the RPC interface can work
			if err != nil {                                                                                 //TODO: Make the error mean something so we can check it here!
				log.Printf("[DIAGNOSTIC] Stabilization call failed!")
				log.Printf("[DIAGNOSTIC] Error: " + strconv.Itoa(int(remote.Key)))
				log.Print(err)
				log.Printf("[DIAGNOSTIC] Assuming that the error is the result of a successor node disconnection. Jumping new successor: " + fingers[1].Address)
				successor = fingers[1]
			}
			remote = remoteInterf.(RemoteNode)
		}
		if remote.Key.BetweenExclusive(my.key, successor.Key) {
			log.Printf("New successor %d\n", remote.Key)
			successor = &remote
			fingers[0] = &remote
			log.Printf("My keyspace is (%d, %d)\n", my.key, successor.Key)
		}

		rpcNotify(joinAddrPort(successor.Address, config.CalleePort()), RemoteNode{
			Address: Address,
			Key:     my.key,
		})
		time.Sleep(time.Second * 1)
	}
}

//FixFingers is the finger-table updater.
//Again, this is a goroutine and never terminates.
func FixFingers() {
	log.Printf("Starting to finger nodes...\n") //hehehe
	currentFingerIndex := uint64(0)
	for true {
		currentFingerIndex++
		currentFingerIndex %= config.NumFingers()
		offset := uint64(math.Pow(2, float64(currentFingerIndex)))
		val := (uint64(my.key) + offset) % config.MaxKey()
		newFinger := FindSuccessor(Key(val))
		//log.Printf("Updating finger %d (pointing to key %d) of %d to point to node %s\n", currentFingerIndex, val, len(Fingers), newFinger.Address)
		if newFinger.Address != fingers[currentFingerIndex].Address {
			log.Printf("Updating finger %d (key %d) of %d to point to node %s (key %d)\n", currentFingerIndex, val, len(fingers)-1, newFinger.Address, newFinger.Key)
		}
		fingers[currentFingerIndex] = &newFinger
		time.Sleep(time.Second * 1)
	}
}

// IsAlive is a heartbeat check.  If this fails, the RPC call will err out.
func IsAlive(void bool) bool {
	return void
}

// CheckPredecessor is a goroutine that keeps tabs on the predecessor and updates itself if the predecessor leaves the network.
func CheckPredecessor() {
	for true {
		if predecessor != nil {
			_, err := rpcIsAlive(joinAddrPort(predecessor.Address, config.CalleePort()), true)
			if err != nil {
				log.Printf("Predecessor " + predecessor.Address + " failed a health check!  Attempting to adjust...")
				log.Print(err)
				predecessor = nil
			}
		}
		time.Sleep(time.Second * 1)
	}
}

// CreateLocalNode creates a local node on its own ring.  It can be inserted into another ring later.
func CreateLocalNode() {
	// Initialize the internal table
	internalTable = NewTable(config.MaxKey())

	// Set the variables of this node.
	var err error
	rpcCaller, err = rpc.NewCaller(config.CallerPort())
	if err != nil {
		panic("rpcCaller failed to initialize")
	}
	rpcCallee, err = rpc.NewCallee(config.CalleePort())
	if err != nil {
		panic("rpcCallee failed to initialize")
	}

	Address = config.Addr()

	my.key = Hash(Address, config.MaxKey())
	log.Printf("Keyspace position %d was derived from IP%s\n", my.key, config.Addr())

	predecessor = nil
	successor = &RemoteNode{
		Address: Address,
		Key:     my.key,
	}
	// Initialize the finger table for the solo ring configuration
	fingers = make([]*RemoteNode, config.NumFingers())
	log.Printf("Finger table size %d was derived from the keyspace size\n", config.NumFingers())
	for i := uint64(0); i < config.NumFingers(); i++ {
		fingers[i] = successor
	}

	// Define all of the RPC functions.
	// For more info, look at Yuchen's caller.go and example_test.go
	// Go's type "system" is going to make me kill myself.
	rpcNotify = rpcCaller.Declare(RemoteNode{}, 0, 1*time.Second)
	rpcFindSuccessor = rpcCaller.Declare(Key(1), RemoteNode{}, 1*time.Second)
	rpcGetPredecessor = rpcCaller.Declare(0, RemoteNode{}, 1*time.Second)
	rpcIsAlive = rpcCaller.Declare(true, true, 1*time.Second)
	rpcGetKey = rpcCaller.Declare("", GetKeyResponse{}, 5*time.Second)
	rpcPutKey = rpcCaller.Declare(PutKeyRequest{}, true, 5*time.Second)
	rpcPutKeyBackup = rpcCaller.Declare(PutKeyRequest{}, 0, 5*time.Second)
	rpcGetKeyRange = rpcCaller.Declare(GetKeyRangeRequest{}, []HashEntry{}, 100*time.Second)

	// Hook the rpcCallee into this node's functions
	rpcCallee.Implement(FindSuccessor)
	rpcCallee.Implement(Notify)
	rpcCallee.Implement(GetPredecessor)
	rpcCallee.Implement(IsAlive)
	rpcCallee.Implement(PutKeyBackup)
	rpcCallee.Implement(GetKey)
	rpcCallee.Implement(PutKey)
	rpcCallee.Implement(GetKeyRange)
}

//GetPredecessor is a getter for the predecessor, implemented for the sake of RPC calls.
//Note that the RPC calling interface does not allow argument-free functions, so this takes
//a worthless int as argument.
func GetPredecessor(void int) RemoteNode {
	//log.Printf("RPC Call to GetPredecessor!\n")
	if predecessor == nil {
		//log.Printf("Returned self node, no predecessor set.\n")
		return RemoteNode{
			Address: Address,
			Key:     my.key,
		}
	}
	//log.Printf("Returned predecessor.\n")
	return *predecessor
}

// Join a ring given a node IP address.
func Join(ring string) {
	log.Printf("Connecting node to network at %s\n", config.Introducer())
	ringCallee := joinAddrPort(ring, config.CalleePort())
	ringSuccessorInterf, err := rpcFindSuccessor(ringCallee, my.key)
	if err != nil {
		log.Printf("[DIAGNOSTIC] Join failed.  Target: %s", ringCallee)
		log.Print(err)
		panic("rpcFindSuccessor failed!")
	}
	ringSuccessor := ringSuccessorInterf.(RemoteNode)
	successor = &ringSuccessor
	fingers[0] = &ringSuccessor
	log.Printf("New successor %d!\n", successor.Key)
	log.Printf("My keyspace is (%d, %d)\n", my.key, successor.Key)
}

func GetKey(keyString string) GetKeyResponse {
	log.Printf("GetKey(%s)\n", keyString)
	if !isLocalResponsible(Hash(keyString, config.MaxKey())) {
		log.Printf("GetKey(%s): sorry, it's none of my business\n", keyString)
		return GetKeyResponse{[]byte{0}, false}
	}
	rv, err := internalTable.Get(keyString)
	if err != nil {
		log.Printf("GetKey(%s): no such key\n", keyString)
		return GetKeyResponse{[]byte{0}, false}
	}
	log.Printf("GetKey(%s): success\n", keyString)
	return GetKeyResponse{rv, true}
}

func PutKey(pkr PutKeyRequest) bool {
	keyString := pkr.KeyString
	data := pkr.Data
	log.Printf("PutKey(%s)\n", keyString)
	if !isLocalResponsible(Hash(keyString, config.MaxKey())) {
		log.Printf("PutKey(%s): sorry, it's none of my business\n", keyString)
		return false
	}
	internalTable.Put(keyString, data)
	log.Printf("PutKey(%s): success\n", keyString)
	return true
}

func PutKeyBackup(pkr PutKeyRequest) int {
	keyString := pkr.KeyString
	data := pkr.Data
	internalTable.Put(keyString, data)
	return 1
}

func GetKeyRange(gkr GetKeyRangeRequest) []HashEntry {
	return internalTable.GetRange(gkr.Start, gkr.End)
}

func Start() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Creating local node @IP %s on its own ring of size %d...\n", config.Addr(), config.MaxKey())
	CreateLocalNode()
	rpcCallee.Start()
	rpcCaller.Start()
	if !config.Creator() {
		Join(config.Introducer())
	}
	log.Printf("Beginning stabilizer...\n")
	go Stabilize()
	go FixFingers()
	go CheckPredecessor()
}

func joinAddrPort(addr string, port uint16) string {
	return fmt.Sprintf("%s:%d", addr, port)
}

func isLocalResponsible(k Key) bool {
	if predecessor == nil {
		return false
	}
	return k.BetweenEndInclusive(predecessor.Key, my.key)
}
