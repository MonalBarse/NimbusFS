package p2p

import (
	"fmt"
	"net"
	"sync"
)

// ----------------------------- Core Structures ----------------------------- //

type TCPPeer struct {
	net.Conn
	outbound bool
	wg       *sync.WaitGroup
}

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener    // listener is a server that listens for incoming connections
	rpcCh    chan RPC        // rpcCh is a channel that will be used to receive messages from the network
	mu       sync.RWMutex    // mu is a mutex that will be used to synchronize access to
	peers    map[string]Peer // peers is a map that will store the peers
}

type TCPTransportOptions struct {
	ListenAddress string
	HandshakeFunc HandshakeFunc    // in this project we are using NOPHandshakeFunc does nothing. But if we want to implement the handshake we can implement it by creating a function and passing it here
	Decoder       Decoder          // Decoder is an interface that defines the methods that a decoder must implement
	OnPeer        func(Peer) error // When new peer is connected, this function does something - here we are doing nothing
}

// ------------------------------- xxxxxxx ----------------------------------- //

// ------------------- Peer and Transport Initialization --------------------- //

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

func NewTCPTransport(options TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: options,
		rpcCh:               make(chan RPC, 1024),
		peers:               make(map[string]Peer),
	}
}

// ------------------------------- xxxxxxx ----------------------------------- //
// ------------------ Methods of TCPPeer for Peer Operations ----------------- //

/* Index
1. CloseStream: Close the stream
2. Send: Send data to the peer
*/

// 1. CloseStream ---------------------------//
func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

// 2. Send ---------------------------//
func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Write(b)
	return err
}

// ------------------------------- xxxxxxx ----------------------------------- //
// ------------- Methods of TCPTransport for Transport Operations ------------ //

/* Index
1. Addr: Get the listening address
2. Consume: a .Consume will return a channel that will be used to receive messages from the network.
3. Close: Close the transport
4. Dial: Dial a remote peer
5. ListenAndAccept: Start listening and accepting connections
*/

// 1. Addr ---------------------------//
func (t *TCPTransport) Addr() string {
	return t.ListenAddress
}

// 2. Consume ---------------------------//
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

// 3. Close ---------------------------//
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// 4. ListenAndAccept ---------------------------//
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	return nil
}

// 5. Dial ---------------------------//
func (t *TCPTransport) Dial(address string) error {
	// Dial connects to the address on the named network.
	conn, err := net.Dial("tcp", address) //	Dial("tcp", "198.51.100.1:80") //	Dial("udp", "[2001:db8::1]:domain") //	Dial("tcp", ":80")
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}

// ------------------------------- xxxxxxx ----------------------------------- //

// -------------------- Internal Methods of TCPTransport --------------------- //

/* Index
1. startAcceptLoop: Start accepting connections
2. handleConn: Handle a new connection
*/

// 1. startAcceptLoop ---------------------------//
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept() // conn is the connection object and err is the error object
		if err != nil {
			fmt.Printf("TCPTransport: failed to accept connection: %v\n", err)
			continue
		}
		go t.handleConn(conn, false) // here after accepting the connection (inbound) we are handling the connection by calling handleConn (outbound = false)
	}
}

// 2. handleConn ---------------------------//
func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	peer := NewTCPPeer(conn, outbound) // creates a new TCPPeer object using the connection

	if err := t.HandshakeFunc(peer); err != nil { //performs a handshake using the HandshakeFunc, here we are using NOPHandshakeFunc
		fmt.Printf("TCPTransport: handshake failed: %v\n", err)
		conn.Close()
		return
	}

	if t.OnPeer != nil { // The OnPeer function is called when a new peer is connected. It is used to do something when a new peer is connected. Here we are doing nothing
		if err := t.OnPeer(peer); err != nil {
			fmt.Printf("TCPTransport: OnPeer failed: %v\n", err)
			conn.Close()
			return
		}
	}

	// The peer is added to the peers map
	t.mu.Lock()
	t.peers[conn.RemoteAddr().String()] = peer // conn.RemoteAddr() will give the address of the remote peer
	t.mu.Unlock()

	for {
		rpc := RPC{}
		err := t.Decoder.Decode(conn, &rpc) //It uses the Decoder to decode the incoming message into the RPC object.
		if err != nil {
			break
		}
		// It sets the From field of the RPC to the remote address. It sends the RPC object to the rpcCh channel for processing.
		rpc.From = conn.RemoteAddr().String()
		t.rpcCh <- rpc
	}

	// If there is an error while decoding the message, the connection is closed and the peer is removed from the peers map.
	t.mu.Lock()
	delete(t.peers, conn.RemoteAddr().String())
	t.mu.Unlock()
	conn.Close()
}

// ------------------------------- xxxxxxx ----------------------------------- //
