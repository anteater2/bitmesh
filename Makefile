.PHONY: all build first node dht stop

all: build

build:
	docker build -t bitmesh -f test/Dockerfile .

first:
	docker run -it bitmesh chord -n 10

node:
	docker run -it bitmesh chord -n 10 -c 172.17.0.2:2001

dht:
	docker run -it bitmesh dht

caller:
	docker run -it bitmesh node_caller

stop:
	docker stop $(shell docker ps -aq)

clean:
	docker rm $(shell docker ps -qa --no-trunc --filter "status=exited")

nc:
	(docker run -t bitmesh chord -n 10 &\
	docker run -t bitmesh chord -n 10 -c 172.17.0.2:2001 &\
	docker run -t bitmesh chord -n 10 -c 172.17.0.2:2001 &\
	docker run -t bitmesh chord -n 10 -c 172.17.0.2:2001 &\
	docker run -t bitmesh chord -n 10 -c 172.17.0.2:2001 &\
	(sleep 30 && docker run -t bitmesh dht))|nc localhost 3000

deepclean:
	docker kill $(shell docker ps -a -q);	docker rm $(shell docker ps -a -q)