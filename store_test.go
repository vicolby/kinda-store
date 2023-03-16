package main

import (
	"bytes"
	"testing"
)

func TestStire(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some data"))
	if err := s.writeStream("somekey", data); err != nil {
		t.Error(err)
	}
}
