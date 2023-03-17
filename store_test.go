package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpicture"
	pathKey := CASPathTransformFunc(key)
	expectedPath := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	expectedFilename := "6804429f74181a63c50c3d81d733a12f14a353ff"
	assert.Equal(t, expectedPath, pathKey.PathName)
	assert.Equal(t, expectedFilename, pathKey.Filename)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	key := "somekey"
	dataBytes := []byte("some data")
	data := bytes.NewReader(dataBytes)

	if err := s.writeStream(key, data); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, dataBytes, b)

}
