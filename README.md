# Bitmesh
Bitmesh is a library to build distributed applications.

## Structure
* [message](./message): Go object transmission.
* [rpc](./rpc): A RPC library
* [chord](./chord): The chord algorithm and interfaces to it.
* [test](./test): Tests that need to run on docker.

## IP Resolution
Each node needs to know its own IP, because the RPC library definitely doesn't.
This is handled in a really hacky and sort of disgusting way: we attempt to connect to Google's DNS servers at 8.8.8.8.
We don't actually care about the DNS server, but when the connection is made we can snoop to see what local IP is bound to it.

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
