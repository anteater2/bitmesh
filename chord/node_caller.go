package chord

import "github.com/anteater2/bitmesh/rpc"
import "time"

// NodeCaller wraps all the rpc call to a ndoe
type NodeCaller struct {
	caller         *rpc.Caller
	isAlive        rpc.RemoteFunc
	notify         rpc.RemoteFunc
	findSuccessor  rpc.RemoteFunc
	getPredecessor rpc.RemoteFunc
	getSuccessor   rpc.RemoteFunc
	getKeyRange    rpc.RemoteFunc
	getFingers     rpc.RemoteFunc
	get            rpc.RemoteFunc
	put            rpc.RemoteFunc
}

// NewNodeCaller creates a new NodeCaller
func NewNodeCaller(port uint16) (*NodeCaller, error) {
	caller, err := rpc.NewCaller(port)
	if err != nil {
		return nil, err
	}
	return &NodeCaller{
		caller:         caller,
		isAlive:        caller.Declare(isAliveCall{}, isAliveReply{}, 1*time.Second),
		notify:         caller.Declare(notifyCall{}, notifyReply{}, 1*time.Second),
		findSuccessor:  caller.Declare(findSuccessorCall{}, findSuccessorReply{}, 1*time.Second),
		getPredecessor: caller.Declare(getPredecessorCall{}, getPredecessorReply{}, 1*time.Second),
		getSuccessor:   caller.Declare(getSuccessorCall{}, getSuccessorReply{}, 1*time.Second),
		getKeyRange:    caller.Declare(getKeyRangeCall{}, getKeyRangeReply{}, 5*time.Second),
		getFingers:     caller.Declare(getFingersCall{}, getFingersReply{}, 1*time.Second),
		get:            caller.Declare(getCall{}, getReply{}, 5*time.Second),
		put:            caller.Declare(putCall{}, putReply{}, 5*time.Second),
	}, nil
}

// Start starts the NodeCaller
func (nc *NodeCaller) Start() {
	nc.caller.Start()
}

// Notice:
// 1. All the functions below are rpc and thus very slow!
// 2. Target node is represented as an address string of form "<IP>:<port>"

// IsAlive check whether the node is alive or not.
func (nc *NodeCaller) IsAlive(node string) bool {
	_, err := nc.isAlive(node, isAliveCall{})
	if err != nil {
		return false
	}
	return true
}

// Notify ...
func (nc *NodeCaller) Notify(node string, remoteNode RemoteNode) error {
	_, err := nc.notify(node, notifyCall{remoteNode})
	if err != nil {
		return err
	}
	return nil
}

// FindSuccessor ...
func (nc *NodeCaller) FindSuccessor(node string, key Key) (RemoteNode, error) {
	reply, err := nc.findSuccessor(node, findSuccessorCall{key})
	if err != nil {
		return RemoteNode{}, err
	}
	return reply.(findSuccessorReply).Node, nil
}

// GetPredecessor ...
func (nc *NodeCaller) GetPredecessor(node string) (RemoteNode, error) {
	reply, err := nc.getPredecessor(node, getPredecessorCall{})
	if err != nil {
		return RemoteNode{}, err
	}
	return reply.(getPredecessorReply).Node, nil
}

// GetSuccessor ...
func (nc *NodeCaller) GetSuccessor(node string) (RemoteNode, error) {
	reply, err := nc.getSuccessor(node, getSuccessorCall{})
	if err != nil {
		return RemoteNode{}, err
	}
	return reply.(getSuccessorReply).Node, nil
}

// GetKeyRange ...
func (nc *NodeCaller) GetKeyRange(node string, start Key, end Key) ([]HashEntry, error) {
	reply, err := nc.getKeyRange(node, getKeyRangeCall{start, end})
	if err != nil {
		return nil, err
	}
	return reply.(getKeyRangeReply).Data, nil
}

// Get ...
func (nc *NodeCaller) Get(node string, k string) ([]byte, error) {
	reply, err := nc.get(node, getCall{k})
	if err != nil {
		return []byte{0}, err
	}
	return reply.(getReply).Value, reply.(getReply).Error
}

// Put ...
func (nc *NodeCaller) Put(node string, k string, v []byte) error {
	reply, err := nc.put(node, putCall{k, v})
	if err != nil {
		return err
	}
	return reply.(putReply).Error
}

// GetFingers ...
func (nc *NodeCaller) GetFingers(node string) ([]RemoteNode, error) {
	reply, err := nc.getFingers(node, getFingersCall{})
	if err != nil {
		return nil, err
	}
	return reply.(getFingersReply).Fingers, nil
}
