package rpc

// Call represents a remote call
type call struct {
	ID         uint32
	CallerAddr string
	Arg        interface{}
}

// Reply represents a reply to a remote call
type reply struct {
	ID  uint32
	Ret interface{}
}
