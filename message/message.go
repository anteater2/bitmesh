package message

import (
	"encoding/gob"
	"fmt"
	"net"
	"reflect"
	"sync"
)

// Send encodes the message using gob and sends it to the address ip:port
func Send(ip string, port int, message interface{}) error {
	remoteAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, remoteAddr)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(conn)
	err = enc.Encode(&message)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

// Receiver is bound to a local address (or more precisely, port number)
// and contains handlers for a set of types.
type Receiver struct {
	localAddr *net.TCPAddr
	handlers  sync.Map
	quit      chan bool
	wg        *sync.WaitGroup
}

// NewReceiver creates a new instance of Receiver
func NewReceiver(port int) (*Receiver, error) {
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	return &Receiver{localAddr: laddr}, nil
}

// Register records a message type T, identified by a value of type T,
// and a corresponding handler of type func(*T). This handler will be called
// when a new message of type T is received.
// Note: this function is not thread-safe, and it should not be called
// while there are receivers running.
func (r *Receiver) Register(message interface{}, handler interface{}) {
	messageType := reflect.TypeOf(message)
	handlerType := reflect.TypeOf(handler)
	if !validateHandlerType(messageType, handlerType) {
		panic(fmt.Sprintf("message: invalid handler type %s for message type %s",
			handlerType, messageType))
	}
	gob.Register(message)
	r.handlers.Store(messageType, handler)
}

// Start starts a go routine that listens to incoming messages
// and dispatches them to their registered handlers.
func (r *Receiver) Start() error {
	listener, err := net.ListenTCP("tcp", r.localAddr)
	if err != nil {
		return err
	}
	r.quit = make(chan bool, 1)
	r.wg = new(sync.WaitGroup)
	r.wg.Add(1)
	// start a go routine to listen to connection
	go func() {
		defer listener.Close()
		defer r.wg.Done()
		newConn := make(chan (net.Conn))
		// start a go routine to put new connections into channel newConn
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					return
				}
				newConn <- conn
			}
		}()
		for {
			select {
			case conn := <-newConn:
				go r.handleConnection(conn)
			case <-r.quit:
				return
			}
		}
	}()
	return nil
}

// Stop signals the Receiver to stop and waits until it really stops
func (r *Receiver) Stop() {
	if r.quit != nil {
		r.quit <- true
		r.quit = nil
		r.wg.Wait()
		r.wg = nil
	}
}

func (r *Receiver) handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		dec := gob.NewDecoder(conn)
		var msg interface{}
		err := dec.Decode(&msg)
		if err != nil {
			return
		}
		go r.handleMessage(msg)
	}
}

func (r *Receiver) handleMessage(message interface{}) {
	messageType := reflect.TypeOf(message)
	handler, ok := r.handlers.Load(messageType)
	if ok {
		handlerValue := reflect.ValueOf(handler)
		handlerValue.Call([]reflect.Value{reflect.ValueOf(message)})
	}
}

func validateHandlerType(messageType reflect.Type, handlerType reflect.Type) bool {
	if handlerType.Kind() != reflect.Func {
		return false
	}
	if handlerType.NumIn() != 1 || handlerType.NumOut() != 0 {
		return false
	}
	if !messageType.AssignableTo(handlerType.In(0)) {
		return false
	}
	return true
}
