package main

import (
	"fmt"
	"log"

	"github.com/MonalBarse/NimbusFS/p2p"
)

func main() {
	fmt.Println("Hello, Nimbus")
	tr := p2p.NewTCPTransport(":3000")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("failed to listen and accept: %v", err)
	}

	select {}
}
