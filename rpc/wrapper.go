package rpc

// Call represents a remote call
type call struct {
	ID           uint32
	Arg          interface{}
	CallerPort   int
	CallerAddr   string
	IsPassedCall bool // indicates whether CallerAddr or sender's address should be used
}

// Reply represents a reply to a remote call
type reply struct {
	ID  uint32
	Ret interface{}
}
