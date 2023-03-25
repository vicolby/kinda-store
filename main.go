package main

import (
	"log"
	"time"

	"github.com/vicolby/kinda-store/p2p"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	fileServer := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(3 * time.Second)
		fileServer.Stop()
	}()

	if err := fileServer.Start(); err != nil {
		log.Fatal(err)
	}
}
