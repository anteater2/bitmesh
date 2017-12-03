package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/anteater2/bitmesh/dht"
)

func main() {
	dht.Start()
	test := false
	for _, a := range flag.Args() {
		if a == "test" {
			test = true
		}
	}
	if test {
		time.Sleep(10 * time.Second)
		max := 1000

		for i := 0; i < max; i++ {
			k := fmt.Sprintf("%d", i)
			err := dht.Put(k, k)
			if err != nil {
				panic(err)
			}
		}
		for i := 0; i < 100; i++ {
			k := fmt.Sprintf("%d", rand.Int()%max)
			v, err := dht.Get(k)
			if err != nil {
				panic(err)
			}
			if string(k) != v {
				panic("unmatch")
			}
		}
		fmt.Println("done")
	}
	select {}
}
