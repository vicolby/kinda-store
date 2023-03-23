package p2p

// Peer is the interface that represents the remote node
type Peer interface {
	Close() error
}

// Transport is anything thats handles the communication
// between the node in network. This can be of the
// form (TCP, UDP, websockets...)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
