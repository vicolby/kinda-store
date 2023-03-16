package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpicture"
	pathKey := CASPathTransformFunc(key)
	expectedPath := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	expectedFilename := "6804429f74181a63c50c3d81d733a12f14a353ff"
	assert.Equal(t, expectedPath, pathKey.PathName)
	assert.Equal(t, expectedFilename, pathKey.Original)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some data"))
	if err := s.writeStream("somekey", data); err != nil {
		t.Error(err)
	}
}
