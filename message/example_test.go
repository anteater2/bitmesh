package message_test

import (
	"fmt"
	"time"

	"github.com/anteater2/bitmesh/message"
)

type myStruct struct {
	Field1 string
	Field2 int
}

func stringHandler2(str string) {
	fmt.Printf("r2 receives string: %v\n", str)
}

func intHandler1(i int) {
	fmt.Printf("r1 receives int: %v\n", i)
}

func myStructHandler1(s myStruct) {
	fmt.Printf("r1 receives myStruct: %v\n", s)
}

func myStructHandler2(s myStruct) {
	fmt.Printf("r2 receives myStruct: %v\n", s)
}

func Example() {
	r1, _ := message.NewReceiver(8888)
	r1.Register(0, intHandler1)
	r1.Register(myStruct{}, myStructHandler1)

	r2, _ := message.NewReceiver(8889)
	r2.Register("", stringHandler2)
	r2.Register(myStruct{}, myStructHandler2)

	r1.Start()
	r2.Start()

	fmt.Printf("sends r2 string: %v\n", "a string")
	message.Send("", 8889, "a string")

	fmt.Printf("sends r1 int: %v\n", 123)
	message.Send("", 8888, 123)

	fmt.Printf("sends r1 myStruct: %v\n", myStruct{"to r1", 2})
	message.Send("", 8888, myStruct{"to r1", 2})

	// wait 10 ms to ensure messages above are received
	time.Sleep(time.Millisecond * 10)

	fmt.Println("closing r2")
	r2.Stop()

	// the following message will never be received because r2 has been stopped
	fmt.Printf("sends r2 string: %v // won't be received\n", myStruct{"to r2", 1})
	message.Send("", 8889, myStruct{"to r2", 1})

	fmt.Println("closing r1")
	r1.Stop()
	// Unordered output:
	//sends r2 string: a string
	// sends r1 int: 123
	// r2 receives string: a string
	// sends r1 myStruct: {to r1 2}
	// r1 receives int: 123
	// r1 receives myStruct: {to r1 2}
	// closing r2
	// sends r2 string: {to r2 1} // won't be received
	// closing r1
}
