package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
// TCPPeer struct gives us the info about the connection and the direction of the connection
type TCPPeer struct {
	// conn is the connection that is established between the local node and the remote node (between the peers)
	conn net.Conn
	// if we dial and connect to a peer then outbound is true,
	// if a peer connects to us then outbound is false
	outbound bool
	// wg is a wait group that is used to wait for the connection to close
	wg *sync.WaitGroup
}

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener
	rpcCh    chan RPC

	mu    sync.RWMutex
	peers map[string]Peer
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close closes the connection between the local node and the remote node
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOptions struct {
	ListenAddress string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder // Decoder will be used to decode the incoming data into RPC (Remote Procedure Call) which will be used to send data between the peers
	OnPeer        func(Peer) error
}

func NewTCPTransport(options TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: options,
		rpcCh:               make(chan RPC),
	}
}

// Consume returns a channel that will be used to receive messages from the network
// it implements a method from the Transport interface
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	return nil

}
func (t *TCPTransport) startAcceptLoop() {
	// Accept incoming connections
	// Once a connection is accepted, it will be handled by handleConn
	for {
		conn, err := t.listener.Accept()

		if err != nil {
			fmt.Printf("TCPTransport: failed to accept connection: %v\n", err)
			continue
		}

		go t.handleConn(conn, false)
	}

}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	fmt.Printf("Handling connection %+v\n", conn)
	// Read the first byte to determine the protocol
	// If the protocol is not supported, close the connection
	// If t the protocol is supported, pass the connection to the appropriate handler

	var err error
	defer func() {
		fmt.Printf("Closing connection %+v\n", conn)
		conn.Close()
	}()
	peer := NewTCPPeer(conn, outbound)
	// fmt.Printf("Peer %+v\n", peer)

	if err = t.HandshakeFunc(peer); err != nil {
		fmt.Printf("TCPTransport: handshake failed: %v\n", err)
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			fmt.Printf("TCPTransport: OnPeer failed: %v\n", err)
			return
		}
	}

	// Read Loop
	msg := RPC{}
	// Read Loop
	for {
		if err = t.Decoder.Decode(conn, &msg); err != nil {
			fmt.Printf("TCPTransport: failed to decode incoming data: %+v\n", err)
			return
		}
		msg.From = conn.RemoteAddr()

		fmt.Printf("TCPTransport: decoded data: %+v\n", msg)
		t.rpcCh <- msg // We are sending the message to the channel
	}

}
