# message
To send and receive go objects through TCP connections.

## Sender
Sender sends data of a particular set of types.
```
type Sender struct {
	// Has unexported fields.
}

func NewSender() *Sender
func (s *Sender) Register(v interface{})
func (s *Sender) Send(addr string, message interface{}) error
```
Detailed documentations can be found in [source file](./sender.go)

## Receiver
Receiver is bound to a local address (or more precisely, port number) and
contains handlers for a set of types.
```
type Receiver struct {
	// Has unexported fields.
}

func NewReceiver(port uint16, handler func(string, interface{})) (*Receiver, error)
func (r *Receiver) Addr() string
func (r *Receiver) Register(v interface{})
func (r *Receiver) Start() error
func (r *Receiver) Stop()
```
Detailed documentations can be found in [source file](./receiver.go)

## Example
See [example_test.go](./example_test.go)