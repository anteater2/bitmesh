package rpc

// Call represents a remote call
type call struct {
	ID         int64
	CallerAddr string
	Arg        interface{}
}

// Reply represents a reply to a remote call
type reply struct {
	ID  int64
	Ret interface{}
}
