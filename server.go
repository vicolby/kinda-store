package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/vicolby/kinda-store/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         *p2p.TCPTransport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts
	Store *Store

	peerLock sync.RWMutex
	peers    map[string]p2p.Peer

	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		Store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

type Message struct {
	From string
	Payload any
}

type DataMessage struct {
	Key  string
	Data []byte
}

func (f *FileServer) broadcast(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range f.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)

	return gob.NewEncoder(mw).Encode(msg)
}

func (f *FileServer) StoreData(key string, r io.Reader) error {

	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)

	if err := f.Store.Write(key, tee); err != nil {
		return err
	}

	p := &DataMessage{
		Key:  key,
		Data: buf.Bytes(),
	}

	return f.broadcast(&Message{
		From: "TODO",
		Payload: p,
	})
}

func (f *FileServer) Stop() {
	close(f.quitch)
}

func (f *FileServer) OnPeer(peer p2p.Peer) error {
	f.peerLock.Lock()
	defer f.peerLock.Unlock()
	f.peers[peer.RemoteAddr().String()] = peer

	log.Printf("New peer: %s", peer.RemoteAddr().String())

	return nil
}

func (f *FileServer) loop() {
	defer func() {
		log.Println("Closing file server")
		f.Transport.Close()
	}()

	for {
		select {
		case <-f.quitch:
			return
		case msg := <-f.Transport.Consumer():
			fmt.Println(msg)
		}
	}
}

func (f *FileServer) 

func (f *FileServer) bootstrapNetwork() error {
	for _, addr := range f.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			fmt.Println("Dialing", addr)
			if err := f.Transport.Dial(addr); err != nil {
				log.Printf("Error dialing %s: %s", addr, err)
			}
		}(addr)
	}
	return nil
}

func (f *FileServer) Start() error {
	if err := f.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if err := f.bootstrapNetwork(); err != nil {
		return err
	}

	f.loop()
	return nil
}
