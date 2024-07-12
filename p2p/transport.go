package p2p

import "net"

/*
Peer is an interface that defines the methods that a peer must implement ( It represents a node in the network )
*/
type Peer interface {
	Close() error      // Close closes the connection between the local node and the remote node
	Send([]byte) error // Send sends a message to the remote node
	net.Conn           // Conn returns the connection between the local node and the remote node
	// All of these merthods are implemented in the TCPPeer struct in tcp_transport.go
}

/*
Transport is an interface that defines the methods that a transport must implement
Transport is anything that handles communication between the nodes in the network (peers)
eg. TCP, UDP, Websockets, etc.
*/
type Transport interface {
	Addr() string           // Addr returns the address of the local node
	Dial(string) error      // Dial dials a remote node and returns a peer
	ListenAndAccept() error // ListenAndAccept listens for incoming connections and accepts them if they are of the correct protocol may it be TCP, UDP  websockets etc
	Close() error           // Close closes the connection between the local node and the remote node
	Consume() <-chan RPC    // Consume returns a channel that will be used to receive messages from the network
	// All of these merthods are implemented in the TCPTransport struct in tcp_transport.go
}
