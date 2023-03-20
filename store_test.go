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
	s := newStore()
	defer tearDown(t, s)

	key := "somekey"
	dataBytes := []byte("some data")
	data := bytes.NewReader(dataBytes)

	if err := s.writeStream(key, data); err != nil {
		t.Error(err)
	}

	assert.True(t, s.Has(key))

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, dataBytes, b)

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

	assert.False(t, s.Has(key))
}

func newStore() *Store {
	opts := StoreOpts{
		Root:              "network",
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
