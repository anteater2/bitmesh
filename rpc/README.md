# RPC
A RPC library supporting indirect return.
It is built based on package message.

## Caller
Caller represents a caller service where remote functions are declared. It
sends the call to callee over the network and captures the correpsponding
return value.
```
type Caller struct {
	// Has unexported fields.
}

func NewCaller(port uint16) (*Caller, error)
func (c *Caller) Declare(arg interface{}, ret interface{}, timeout time.Duration) RemoteFunc
func (c *Caller) Start() error
func (c *Caller) Stop()
```
Detailed documentations can be found in [source file](./caller.go).

## Callee
Callee represents a callee service where remote functions are implemented.
```
type Callee struct {
	// Has unexported fields.
}    

func NewCallee(port uint16) (*Callee, error)
func (c *Callee) Implement(f interface{})
func (c *Callee) Start() error
func (c *Callee) Stop()
```
Detailed documentations can be found in [source file](./callee.go).

## Example
See [example_test.go](./example_test.go)