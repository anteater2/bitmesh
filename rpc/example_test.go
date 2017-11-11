package rpc_test

import (
	"fmt"
	"time"

	"github.com/anteater2/bitmesh/rpc"
)

// Declare (on the caller side and callee side)
// Note: only exported fields are transferred
type addArg struct {
	X int
	Y int
}

type mulArg struct {
	X int
	Y int
}

// Implement (only on the callee side)
func doAddition(arg addArg) int {
	return arg.X + arg.Y
}

func doMultiplication(arg mulArg) int {
	return arg.X * arg.Y
}

func Example() {
	// setup caller
	caller, _ := rpc.NewCaller(2000)
	// caller needs to declare remote functions
	// and specify return type and timeout
	add := caller.Declare(addArg{}, 0, time.Second)
	mul := caller.Declare(mulArg{}, 0, time.Second)

	// setup callees
	callee1, _ := rpc.NewCallee(2001)
	callee2, _ := rpc.NewCallee(2002)

	// callee needs to implement remote functions
	callee1.Implement(doAddition)
	// callee1 hands over multiplication to callee2
	callee1.Implement(func(arg mulArg, pass rpc.PassFunc) (int, bool) {
		pass("localhost:2002", arg)
		return 0, false
	})
	callee2.Implement(doMultiplication)

	// start
	caller.Start()
	callee1.Start()
	callee2.Start()

	// the function declared is called with (address string, arg interface{})
	// this will be handled by callee1
	res, err := add("localhost:2001", addArg{1, 2})
	if err != nil {
		panic(err)
	}
	sum := res.(int)
	fmt.Printf("1 + 2 = %d\n", sum)

	// this will be handled by callee2
	res, err = mul("localhost:2001", mulArg{3, 4})
	if err != nil {
		panic(err)
	}
	prod := res.(int)
	fmt.Printf("3 * 4 = %d\n", prod)

	// stop
	caller.Stop()
	callee1.Stop()
	callee2.Stop()

	// Output:
	// 1 + 2 = 3
	// 3 * 4 = 12
}
