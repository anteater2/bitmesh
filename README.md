# bitmesh

## DHT test
### Build
Under the directory bitmesh, run following:
```
docker build -t bitmesh -f script/Dockerfile .
```

### Start first node
```
docker run -it bitmesh bitmesh -n 10
```

### Start more nodes
```
docker run -it bitmesh bitmesh -n 10 -c 172.17.0.2
```

### Start test node
```
docker run -it bitmesh bitmesh -n 10 -c 172.17.0.2 test
```
Wait 10 seconds for this node to setup, then it will try to
put 1000 key value pairs and randomly get 100 out of them.
If any of them are lost or changed, the program will panic.