.PHONY: all build first node dht stop

all: build

build:
	docker build -t bitmesh -f test/Dockerfile .

first:
	docker run -it bitmesh chord -n 10

node:
	docker run -it bitmesh chord -n 10 -c 172.17.0.2

dht:
	docker run -it bitmesh dht

stop:
	docker stop $(shell docker ps -aq)

clean:
	docker rm $(shell docker ps -qa --no-trunc --filter "status=exited")