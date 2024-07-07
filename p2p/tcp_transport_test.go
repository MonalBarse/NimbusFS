package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	tcpOptions := TCPTransportOptions{
		ListenAddress: ":8080",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(tcpOptions)
	assert.Equal(t, ":8080", tr.ListenAddress)
	assert.Nil(t, tr.ListenAndAccept())

	// Server
	// tr.Start()
}
