# Bitmesh DHT

## Tests on Docker network
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

To run dht test,
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

To run a full network test and send the output through a network socket:
```
make nc
```
This runs first, node, and dht and pipes the aggregate output to localhost:3000 (you can see this by running ```nc -l -p 3000```). Note that if there is nothing listening there, this won't work.

This waits 30 seconds before starting dht so that the finger tables have time to fully initialize.

If this broke and you need to kill a lot of docker containers, try:
```
mk deepclean
```
*THIS WILL KILL AND REMOVE EVERY DOCKER CONTAINER ON YOUR MACHINE*