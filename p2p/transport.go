package p2p

import "net"

// Peer is the interface that represents the remote node
type Peer interface {
	net.Conn
	Send([]byte) error
}

// Transport is anything thats handles the communication
// between the node in network. This can be of the
// form (TCP, UDP, websockets...)
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
