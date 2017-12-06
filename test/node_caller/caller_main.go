package main

import (
	"fmt"

	"github.com/anteater2/bitmesh/chord"
)

func main() {
	caller, err := chord.NewNodeCaller(2000)
	if err != nil {
		panic(err)
	}
	caller.Start()
	node := "172.17.0.2:2001"
	fingers, err := caller.GetFingers(node)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got fingers from %s (%v):\n", node, chord.Hash(node, 1<<10))
	for i, n := range fingers {
		fmt.Printf("%d: %v (%v)\n", i, n.Address, n.Key)
	}
}
