# Chord Node

The algorithem is implemented according to MIT Chord Paper::
https://pdos.csail.mit.edu/papers/ton:chord/paper-ton.pdf

[The RPC library](../rpc/) is used.

# The Chord algorithm
See [node.go](./node.go)
### Interface.
```
func Start(addr string, calleePort uint16, callerPort uint16, bits uint64) error
func Join(ring string) error
```

### Ports
Callers send on port 2000.
Callees receive on port 2001.

### Fault tolerance
A fully calibrated/set up ring should be able to handle a single node going offline without losing data or breaking.<br>
This doesn't mean that nodes can be removed frequently; if a node fails, the network has to fix its successor lists and otherwise adjust before it can tolerate another one.

# Client
NodeCaller wraps all the rpc call to a ndoe.
```
type NodeCaller struct {
	// Has unexported fields.
}

func NewNodeCaller(port uint16) (*NodeCaller, error)
func (nc *NodeCaller) FindSuccessor(node string, key Key) (RemoteNode, error)
func (nc *NodeCaller) Get(node string, k string) ([]byte, error)
func (nc *NodeCaller) GetFingers(node string) ([]RemoteNode, error)
func (nc *NodeCaller) GetKeyRange(node string, start Key, end Key) ([]HashEntry, error)
func (nc *NodeCaller) GetPredecessor(node string) (RemoteNode, error)
func (nc *NodeCaller) GetSuccessor(node string) (RemoteNode, error)
func (nc *NodeCaller) IsAlive(node string) bool
func (nc *NodeCaller) Notify(node string, remoteNode RemoteNode) error
func (nc *NodeCaller) Put(node string, k string, v []byte) error
func (nc *NodeCaller) Start()
```
See [node_caller.go](./node_caller.go)