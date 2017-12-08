package chord

import (
	"log"

	"github.com/anteater2/bitmesh/rpc"
)

var callee *rpc.Callee

func initNodeCallee(port uint16) error {
	var err error
	callee, err = rpc.NewCallee(port)
	if err != nil {
		return err
	}

	callee.Implement(handleIsAliveCall)
	callee.Implement(handleNotifyCall)
	callee.Implement(handleFindSuccessor)
	callee.Implement(handleGetFingers)
	callee.Implement(handleGet)
	callee.Implement(handlePut)
	callee.Implement(handleGetPredecessor)
	callee.Implement(handleGetSuccessor)
	callee.Implement(handleGetKeyRange)

	return nil
}

func startNodeCallee() {
	callee.Start()
}

// ----------------------------------------------------------------------------

type isAliveCall struct{}

type isAliveReply struct{}

func handleIsAliveCall(call isAliveCall) isAliveReply {
	return isAliveReply{}
}

// ----------------------------------------------------------------------------

type notifyCall struct {
	RemoteNode RemoteNode
}

type notifyReply struct{}

func handleNotifyCall(call notifyCall) notifyReply {
	notify(call.RemoteNode)
	return notifyReply{}
}

// ----------------------------------------------------------------------------

type findSuccessorCall struct {
	Key Key
}

type findSuccessorReply struct {
	Node RemoteNode
}

func handleFindSuccessor(call findSuccessorCall, pass rpc.PassFunc) (findSuccessorReply, bool) {
	key := call.Key
	if key.BetweenEndInclusive(my.key, successor.Key) {
		return findSuccessorReply{*successor}, true
	}
	target := closestPrecedingNode(key)
	if target.Address == my.address {
		log.Printf("[DIAGNOSTIC] Infinite loop detected!\n")
		log.Printf("[DIAGNOSTIC] This is likely because of a bad finger table.\n")
		target = *successor
	}
	pass(target.Address, call)
	return findSuccessorReply{}, false
}

// ----------------------------------------------------------------------------

type getCall struct {
	Key string
}

type getReply struct {
	Value []byte
	Error error
}

func handleGet(call getCall) getReply {
	rv, err := getKey(call.Key)
	return getReply{rv, err}
}

// ----------------------------------------------------------------------------

type putCall struct {
	Key   string
	Value []byte
}

type putReply struct {
	Error error
}

func handlePut(call putCall) putReply {
	err := putKey(call.Key, call.Value)
	return putReply{err}
}

// ----------------------------------------------------------------------------

type getPredecessorCall struct{}

type getPredecessorReply struct {
	Node RemoteNode
}

func handleGetPredecessor(call getPredecessorCall) getPredecessorReply {
	return getPredecessorReply{getPredecessor()}
}

// ----------------------------------------------------------------------------

type getSuccessorCall struct{}

type getSuccessorReply struct {
	Node RemoteNode
}

func handleGetSuccessor(call getSuccessorCall) getSuccessorReply {
	return getSuccessorReply{*successor}
}

// ----------------------------------------------------------------------------

type getKeyRangeCall struct {
	Start Key
	End   Key
}

type getKeyRangeReply struct {
	Data []HashEntry
}

func handleGetKeyRange(call getKeyRangeCall) getKeyRangeReply {
	return getKeyRangeReply{getKeyRange(call.Start, call.End)}
}

// ----------------------------------------------------------------------------

type getFingersCall struct{}

type getFingersReply struct {
	Fingers []RemoteNode
}

func handleGetFingers(call getFingersCall) getFingersReply {
	reply := getFingersReply{}
	reply.Fingers = make([]RemoteNode, len(fingers))
	for i, e := range fingers {
		reply.Fingers[i] = *e
	}
	return reply
}
