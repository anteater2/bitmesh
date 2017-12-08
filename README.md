# Bitmesh Chord Node
## Structure
```node.go``` is the only file that is responsible for everything.
## Protocol
Node-to-node communications use RPC based around the Go object serialization format (gob).

Some RPC calls are documented in the MIT Chord Paper:

https://pdos.csail.mit.edu/papers/ton:chord/paper-ton.pdf

Extra RPC calls exist to allow for actual Get/Put.
## Implementation
Each node implements callee receivers for:

```FindSuccessor(key uint32)```

```Notify(node RemoteNode)```

```GetPredecessor()```

```IsAlive() //but this one might be hard```

They should also define caller interfaces for each. The periodic functions:

```Stabilize```

```FixFingers```

```CheckPredecessor // again, this is hard, and currently not implemented```

are goroutines with sleeps.

## IP Resolution
Each node needs to know its own IP, because the RPC library definitely doesn't.
This is handled in a really hacky and sort of disgusting way: we attempt to connect to Google's DNS servers at 8.8.8.8.
We don't actually care about the DNS server, but when the connection is made we can snoop to see what local IP is bound to it.

## Ports
Callers send on port 2000.
Callees receive on port 2001.

## Fault tolerance
A fully calibrated/set up ring should be able to handle a single node going offline without losing data or breaking.<br>
This doesn't mean that nodes can be removed frequently; if a node fails, the network has to fix its successor lists and otherwise adjust before it can tolerate another one.

# Docker Testing
See the directory test. Or,

To build docker images,
```
make
```

To run first chord node,
```
make first
```

To run more chord nodes,
```
make node
```

To run DHT test,
```
make dht
```

To run node caller test,
```
make caller
```

To stop all containers,
```
make stop
```

To clean up containers,
```
make clean
```

To run a full DHT test and send the output through a network socket:
```
make nc
```
This runs first, node, and dht and pipes the aggregate output to localhost:3000 (you can see this by running ```nc -l 3000```). Note that if there is nothing listening there, this won't work.

This waits 30 seconds before starting dht so that the finger tables have time to fully initialize.

If this broke and you need to kill a lot of docker containers, try:<br>
**CAUTION: THIS WILL KILL AND REMOVE EVERY DOCKER CONTAINER ON YOUR MACHINE!**
```
mk deepclean
```
