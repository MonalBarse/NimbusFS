package p2p

type HandshakeFunc func(Peer) error

// NOPHandshakeFunc are generally used for testing purposes
// It is a no-op function - which means it does nothing
// use case eg- when we want to test the transport layer without actually doing the handshake
func NOPHandshakeFunc(Peer) error {
	return nil
}
