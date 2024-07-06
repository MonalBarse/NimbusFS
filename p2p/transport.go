package p2p

/*
Peer is an interface that defines the methods that a peer must implement ( It represents a node in the network )
*/
type Peer struct {
}

/*
Transport is an interface that defines the methods that a transport must implement
Transport is anything that handles communication between the nodes in the network (peers)
eg. TCP, UDP, Websockets, etc.
*/
type Transport interface {
	ListenAndAccept() error // ListenAndAccept listens for incoming connections and accepts them if they are of the correct protocol may it be TCP, UDP  websockets etc

}
