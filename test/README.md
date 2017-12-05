# Test on Docker

## Build
Under the directory **bitmesh**, run following:
```
docker build -t bitmesh -f test/Dockerfile .
```

## Run chord node
```
docker run -it bitmesh chord -n 10 [-c 172.17.0.2]
```

## Run dht test (must have some chord nodes running)
```
docker run -it bitmesh dht
```
The client will try to put 1000 key value pairs and randomly get 100 out of them.
If any of them are lost or changed, the program will panic.
