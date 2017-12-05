package main

import (
	"fmt"
	"math/rand"

	"github.com/anteater2/bitmesh/dht"
)

func main() {
	t, err := dht.New("172.17.0.2:2001", 2001, 2000)
	if err != nil {
		panic(err)
	}
	t.Start()

	max := 1000
	for i := 0; i < max; i++ {
		k := fmt.Sprintf("%d", i)
		err := t.Put(k, k)
		if err != nil {
			panic(err)
		}
	}
	for i := 0; i < 100; i++ {
		k := fmt.Sprintf("%d", rand.Int()%max)
		v, err := t.Get(k)
		if err != nil {
			panic(err)
		}
		if string(k) != v {
			panic("unmatch")
		}
	}
	fmt.Println("success")
}
