package message

import (
	"encoding/gob"
	"fmt"
	"net"
	"reflect"
	"sync"
)

var sendable = make(map[reflect.Type]bool)
var mutex sync.Mutex

// Register records a type so that package message can send it
func Register(v interface{}) {
	mutex.Lock()
	gob.Register(v)
	sendable[reflect.TypeOf(v)] = true
	mutex.Unlock()
}

// Send encodes the message using gob and sends it to the address ip:port
func Send(addr string, message interface{}) error {
	if _, prs := sendable[reflect.TypeOf(message)]; !prs {
		return fmt.Errorf("message: unregistered type %T", message)
	}
	remoteAddr, err := net.ResolveTCPAddr("tcp", addr)
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
