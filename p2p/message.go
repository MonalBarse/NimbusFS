package p2p

import "net"

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

// What is RPC? - Remote Procedure Call.
// It is a protocol that one program can use to request a service from a program located in another computer in a network without having to understand the network's details.
// RPC holds any arbitrary data that is being send over the each transport between two nodes.

type RPC struct {
	From net.Addr
	// Size    int64
	// Stream  bool
	Payload []byte
}
