package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type TCPPeer struct {
	// conn is the underlying connection of the peer
	net.Conn
	// if we dial and retrieve a conn -> outbound == true
	// if we accept and retrieve a conn -> outbound == false
	outbound bool

	Wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		Wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandShakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

func (t *TCPTransport) Consumer() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)

	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	log.Printf("Listening on %s", t.ListenAddr)

	return nil
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error %s", err)
		}

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("Closing conn %s \n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

	if err = t.HandShakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("Handshake error %s", err)
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	rpc := RPC{}

	for {
		err = t.Decoder.Decode(conn, &rpc)

		if err != nil {
			fmt.Printf("Decode error %s ", err)
			return
		}

		rpc.From = conn.RemoteAddr()
		peer.Wg.Add(1)
		fmt.Println("Waiting till stream is done")
		t.rpcch <- rpc
		peer.Wg.Wait()
		fmt.Println("Stream is done")
	}
}
