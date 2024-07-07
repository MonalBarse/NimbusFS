package main

import (
	"fmt"
	"log"

	"github.com/MonalBarse/NimbusFS/p2p"
)

func OnPeer(p p2p.Peer) error {
	p.Close()
	fmt.Println("doing somting with peer outside of TCPTransport")
	return nil
}

func main() {
	tcpOptions := p2p.TCPTransportOptions{
		ListenAddress: ":8080",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	fmt.Println("Hello, Nimbus")
	tr := p2p.NewTCPTransport(tcpOptions)

	go func() {
		for {
			rpc := <-tr.Consume()
			fmt.Println("Received message: ", string(rpc.Payload))
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("failed to listen and accept: %+v", err)

	}
	select {}
}
