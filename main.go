package main

import (
	"bytes"
	"log"
	"time"

	"github.com/vicolby/kinda-store/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":3001", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(3 * time.Second)
	go s2.Start()
	time.Sleep(3 * time.Second)

	data := bytes.NewBuffer([]byte("Hello World!"))
	s2.StoreData("hello", data)

	select {}
}
