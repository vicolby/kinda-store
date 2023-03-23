package main

import (
	"fmt"
	"log"

	"github.com/vicolby/kinda-store/p2p"
)

func OnPeer(peer p2p.Peer) error {
	fmt.Println("Doing some logic outside TCPTransport")
	return nil
}

func main() {
	tpcOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tpcOpts)

	go func() {
		for rpc := range tr.Consumer() {
			log.Println(rpc)
		}
	}()
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
