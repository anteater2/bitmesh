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
	fmt.Printf("Exploring node %s (key %v)\n", node, chord.Hash(node, 1<<10))

	pred, err := caller.GetPredecessor(node)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Predecessor: %s (key %v)\n", pred.Address, pred.Key)

	succ, err := caller.GetSuccessor(node)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successor:   %s (key %v)\n", succ.Address, succ.Key)

	fingers, err := caller.GetFingers(node)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Fingers:\n")
	for i, n := range fingers {
		fmt.Printf("[%d]: %v (key %v)\n", i, n.Address, n.Key)
	}
}
