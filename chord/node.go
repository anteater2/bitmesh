package chord

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
)

var internalTable *HashTable

var my struct {
	key     Key
	address string
}

var fingers []*RemoteNode
var predecessor *RemoteNode
var successor *RemoteNode
var caller *NodeCaller

// RemoteNode holds information for connecting to a remote node
type RemoteNode struct {
	Address string
	Key     Key
}

// CreateLocalNode creates a local node on its own ring.  It can be inserted into another ring later.
func createLocalNode() {
	// Initialize the internal table
	internalTable = NewTable(config.maxKey)

	// Set the variables of this node.
	var err error
	caller, err = NewNodeCaller(config.callerPort)
	if err != nil {
		panic("rpcCaller failed to initialize")
	}
	err = initNodeCallee(config.calleePort)
	if err != nil {
		panic("rpcCallee failed to initialize")
	}

	my.address = fmt.Sprintf("%s:%d", config.addr, config.calleePort)

	my.key = Hash(my.address, config.maxKey)
	log.Printf("[NODE %d] Keyspace position %d was derived from address %s\n", my.key, my.key, my.address)

	predecessor = nil
	successor = &RemoteNode{
		Address: my.address,
		Key:     my.key,
	}
	// Initialize the finger table for the solo ring configuration
	fingers = make([]*RemoteNode, config.numFingers)
	log.Printf("[NODE %d] Finger table size %d was derived from the keyspace size\n", my.key, config.numFingers)
	for i := uint64(0); i < config.numFingers; i++ {
		fingers[i] = successor
	}
}

// Start ...
func Start(addr string, calleePort uint16, callerPort uint16, bits uint64) error {
	err := Init(addr, calleePort, callerPort, bits)
	if err != nil {
		return err
	}
	log.Printf("Creating local node @IP %s on its own ring of size %d...\n", config.addr, config.maxKey)
	createLocalNode()
	startNodeCallee()
	caller.Start()
	log.Printf("Beginning stabilizer...\n")
	go stabilize()
	go fixFingers()
	go checkPredecessor()
	return nil
}

// Join a ring given a node IP address.
func Join(ring string) error {
	log.Printf("[NODE %d] Connecting node to network at %s\n", my.key, config.introducer)
	ringSuccessor, err := caller.FindSuccessor(ring, my.key)
	if err != nil {
		return err
	}
	successor = &ringSuccessor
	fingers[0] = &ringSuccessor
	log.Printf("[NODE %d] New successor %d!\n", my.key, successor.Key)
	log.Printf("[NODE %d] My keyspace is (%d, %d)\n", my.key, my.key, successor.Key)
	return nil
}

// closestPrecedingNode finds the closest preceding node to the key in this node's finger table.
// This doesn't need any RPC.
func closestPrecedingNode(key Key) RemoteNode {
	for i := config.numFingers - 1; i > 0; i-- { // WARNING: GO DOES THIS i>0 CHECK AT THE END OF THE LOOP!
		//log.Printf("Checking finger %d\n", i)
		if fingers[i] == nil {
			panic("You attempted to find ClosestPrecedingNode without an initialized finger table!")
		}
		if fingers[i].Key.BetweenExclusive(my.key, key) {
			return *fingers[i]
		}
	}
	return RemoteNode{Address: my.address, Key: my.key}
}

func isLocalResponsible(k Key) bool {
	if predecessor == nil {
		return false
	}
	return k.BetweenEndInclusive(predecessor.Key, my.key)
}

/*****************************************************************************
 * Remote interfaces                                                         *
 *****************************************************************************/

// findSuccessor finds the successor node to the key.  This may require RPC calls.
func findSuccessor(key Key) RemoteNode {
	if key.BetweenEndInclusive(my.key, successor.Key) {
		// key is between this node and its successor
		return *successor
	}
	target := closestPrecedingNode(key)
	if target.Address == my.address {
		log.Printf("[DIAGNOSTIC] Infinite loop detected!\n")
		log.Printf("[DIAGNOSTIC] This is likely because of a bad finger table.\n")
		panic("This is probably a serious bug.")
	}
	// Now, we have to do an RPC on target to find the successor.

	rv, err := caller.FindSuccessor(target.Address, key)
	if err != nil {
		log.Printf("[DIAGNOSTIC] Remote target is " + target.Address + "\n")
		log.Print(err)
		panic("rpcFindSuccessor failed!")
	}
	return rv
}

// get notified
func notify(node RemoteNode) {
	if predecessor == nil || node.Key.BetweenExclusive(predecessor.Key, my.key) {
		log.Printf("[NODE %d] Got notify from %s!  New predecessor: %d\n", my.key, node.Address, node.Key)
		predecessor = &node
		if predecessor.Address != my.address {
			rv, err := caller.GetKeyRange(predecessor.Address, my.key, predecessor.Key)
			if err != nil {
				log.Fatal(err)
			}
			for _, entry := range rv {
				internalTable.Put(entry.Key, entry.Value)
			}
		}
	}
}

// GetPredecessor is a getter for the predecessor, implemented for the sake of RPC calls.
func getPredecessor() RemoteNode {
	//log.Printf("RPC Call to GetPredecessor!\n")
	if predecessor == nil {
		//log.Printf("Returned self node, no predecessor set.\n")
		return RemoteNode{
			Address: my.address,
			Key:     my.key,
		}
	}
	//log.Printf("Returned predecessor.\n")
	return *predecessor
}

// func putKeyBackup(pkr putCall) int {
// 	keyString := pkr.KeyString
// 	data := pkr.Data
// 	internalTable.Put(keyString, data)
// 	return 1
// }

func getKey(keyString string) ([]byte, error) {
	if !isLocalResponsible(Hash(keyString, config.maxKey)) {
		log.Printf("[NODE %d] GetKey(%s): sorry, it's none of my business\n", my.key, keyString)
		return []byte{0}, fmt.Errorf("wrong node to get the key")
	}
	rv, err := internalTable.Get(keyString)
	if err != nil {
		log.Printf("[NODE %d] GetKey(%s): no such key\n", my.key, keyString)
		return []byte{0}, fmt.Errorf("no such key")
	}
	log.Printf("[NODE %d] GetKey(%s): success\n", my.key, keyString)
	return rv, nil
}

func putKey(key string, value []byte) error {
	if !isLocalResponsible(Hash(key, config.maxKey)) {
		log.Printf("[NODE %d] PutKey(%s): sorry, it's none of my business\n", my.key, key)
		return fmt.Errorf("wrong node to get the key")
	}
	internalTable.Put(key, value)
	log.Printf("[NODE %d] PutKey(%s): success\n", my.key, key)
	return nil
}

func getKeyRange(start Key, end Key) []HashEntry {
	return internalTable.GetRange(start, end)
}

/*****************************************************************************
 * Periodically run                                                          *
 *****************************************************************************/

// CheckPredecessor is a goroutine that keeps tabs on the predecessor and updates itself if the predecessor leaves the network.
func checkPredecessor() {
	for true {
		if predecessor != nil {
			if !caller.IsAlive(predecessor.Address) {
				log.Printf("Predecessor " + predecessor.Address + " failed a health check!  Attempting to adjust...")
				predecessor = nil
			}
		}
		time.Sleep(time.Second * 1)
	}
}

// stabilize the Successor and Predecessor fields of this node.
// This is a goroutine and never terminates.
func stabilize() {
	for true { // This is how while loops work.  Not even joking.
		var remote RemoteNode
		var err error
		if predecessor == nil {
			log.Printf("Null predecessor!  New predecessor: %d\n", successor.Key)
			predecessor = successor
		}
		if successor.Address == my.address {
			// Avoid making an RPC call to ourselves
			remote = *predecessor
		} else {
			remote, err = caller.GetPredecessor(successor.Address)
			if err != nil { //TODO: Make the error mean something so we can check it here!
				log.Printf("[DIAGNOSTIC] Stabilization call failed!")
				log.Printf("[DIAGNOSTIC] Error: " + strconv.Itoa(int(remote.Key)))
				log.Print(err)
				log.Printf("[DIAGNOSTIC] Assuming that the error is the result of a successor node disconnection. Jumping new successor: " + fingers[1].Address)
				successor = fingers[1]
			}
		}
		if remote.Key.BetweenExclusive(my.key, successor.Key) {
			log.Printf("New successor %d\n", remote.Key)
			successor = &remote
			fingers[0] = &remote
			log.Printf("My keyspace is (%d, %d)\n", my.key, successor.Key)
		}
		caller.Notify(successor.Address, RemoteNode{
			Address: my.address,
			Key:     my.key,
		})
		time.Sleep(time.Second * 1)
	}
}

// fixFingers is the finger-table updater.
// Again, this is a goroutine and never terminates.
func fixFingers() {
	log.Printf("Starting to finger nodes...\n") //hehehe
	currentFingerIndex := uint64(0)
	for true {
		currentFingerIndex++
		currentFingerIndex %= config.numFingers
		offset := uint64(math.Pow(2, float64(currentFingerIndex)))
		val := (uint64(my.key) + offset) % config.maxKey
		newFinger := findSuccessor(Key(val))
		//log.Printf("Updating finger %d (pointing to key %d) of %d to point to node %s\n", currentFingerIndex, val, len(Fingers), newFinger.Address)
		if newFinger.Address != fingers[currentFingerIndex].Address {
			log.Printf("Updating finger %d (key %d) of %d to point to node %s (key %d)\n", currentFingerIndex, val, len(fingers)-1, newFinger.Address, newFinger.Key)
		}
		fingers[currentFingerIndex] = &newFinger
		time.Sleep(time.Second * 1)
	}
}
