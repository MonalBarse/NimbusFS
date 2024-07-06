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
}

type TCPTransport struct {
	listenAddress string
	listner       net.Listener
	mu            sync.RWMutex
	peer          map[net.Addr]Peer // net.Addr - Network()+String()- for eg tcp://192.0.2.1:25

}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func NewTCPTransport(listenAddress string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddress,
		peer:          make(map[net.Addr]Peer),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listner, err = net.Listen("tcp", t.listenAddress)
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
		conn, err := t.listner.Accept()
		if err != nil {
			fmt.Printf("TCPTransport: failed to accept connection: %v\n", err)
			continue
		}

		go t.handleConn(conn)
	}

}
func (t *TCPTransport) handleConn(conn net.Conn) {
	// Read the first byte to determine the protocol
	// If the protocol is not supported, close the connection
	// If the protocol is supported, pass the connection to the appropriate handler
	fmt.Printf("Handling connection %+v\n", conn)

	peer := NewTCPPeer(conn, true)
	fmt.Printf("Peer %+v\n", peer)

}
