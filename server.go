package main

import (
	"fmt"
	"log"

	"github.com/vicolby/kinda-store/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         *p2p.TCPTransport
}

type FileServer struct {
	FileServerOpts
	Store *Store

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
	}
}

func (f *FileServer) Stop() {
	close(f.quitch)
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

func (f *FileServer) Start() error {
	if err := f.Transport.ListenAndAccept(); err != nil {
		return err
	}

	f.loop()
	return nil
}
