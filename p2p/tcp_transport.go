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
	net.Conn
	// if we dial and connect to a peer then outbound is true,
	// if a peer connects to us then outbound is false
	outbound bool
	// wg is a wait group that is used to wait for the connection to close
	wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener
	rpcCh    chan RPC

	mu    sync.RWMutex
	peers map[string]Peer
}

// Close closes the connection between the local node and the remote node
func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

func (p *TCPPeer) Send(b []byte) error {
	_, erer := p.Conn.Write(b)
	return erer
}

type TCPTransportOptions struct {
	ListenAddress string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder          // Decoder will be used to decode the incoming data into RPC (Remote Procedure Call) which will be used to send data between the peers
	OnPeer        func(Peer) error // OnPeer is a function that will be called when a new peer is connected
}

func NewTCPTransport(options TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: options,
		rpcCh:               make(chan RPC, 1024),
	}
}

// Addr returns the address of the local node
// it implements a method from the Transport interface, the trnasport is accepting connections on this address
func (t *TCPTransport) Addr() string {
	return t.ListenAddress
}

// Consume returns a channel that will be used to receive messages from the network
// it implements a method from the Transport interface
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial dials a remote node and returns a peers
// it implements a method from the Transport interface

func (t *TCPTransport) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}

// ListenAndAccept listens for incoming connections and accepts them if they are of the correct protocol may it be TCP, UDP  websockets etc
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
	for {
		message := RPC{}
		err = t.Decoder.Decode(conn, &message) // Decode the incoming data into RPC (Remote Procedure Call) which will be used to send data between the peers
		if err != nil {
			return
		}
		message.From = conn.RemoteAddr().String()
		if message.Stream {
			peer.wg.Add(1)
			fmt.Printf("Stream message incoming.... %s\n", conn.RemoteAddr())
			peer.wg.Wait()

			fmt.Printf("%s Stream Closed, resuming reading loop\n", conn.RemoteAddr())
			continue

		}
	}

}
