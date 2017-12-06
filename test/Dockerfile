FROM golang:1.9.2

ADD . /go/src/github.com/anteater2/bitmesh

RUN go install github.com/anteater2/bitmesh/test/chord
RUN go install github.com/anteater2/bitmesh/test/dht
RUN go install github.com/anteater2/bitmesh/test/node_caller
