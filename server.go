package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

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
	From    string
	Payload any
}

type MessageStoreFile struct {
	Key  string
	Size int64
}

func (f *FileServer) StoreData(key string, r io.Reader) error {
	var (
		fileBuffer = new(bytes.Buffer)
		tee        = io.TeeReader(r, fileBuffer)
	)

	size, err := f.Store.Write(key, tee)
	if err != nil {
		return err
	}

	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}

	msgBuf := new(bytes.Buffer)

	if err := gob.NewEncoder(msgBuf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range f.peers {
		if err := peer.Send(msgBuf.Bytes()); err != nil {
			return err
		}
	}

	time.Sleep(time.Second * 3)

	for _, peer := range f.peers {
		n, err := io.Copy(peer, fileBuffer)
		if err != nil {
			return err
		}

		fmt.Println("Received and written bytes to disk:", n)
	}

	return nil
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
		case rpc := <-f.Transport.Consumer():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println(err)
			}

			if err := f.handleMessage(rpc.From.String(), &msg); err != nil {
				log.Println(err)
			}
		}
	}
}

func (f *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return f.handleMessageStoreFile(from, v)
	}
	return nil
}

func (f *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	peer, ok := f.peers[from]
	if !ok {
		return fmt.Errorf("could not find (%s) peer in peer list", peer)
	}

	n, err := f.Store.Write(msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		return err
	}

	log.Printf("wrote %d bytes to disk\n", n)

	peer.(*p2p.TCPPeer).Wg.Done()

	return nil
}

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

func init() {
	gob.Register(MessageStoreFile{})
}
