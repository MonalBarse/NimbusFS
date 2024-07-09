// print the data
package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ------------------------ Utility func ------------------------ //
func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func cleanupTest(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
} // for cleanup, what files we have created in the store in the test, we need to remove them

// ------------------------ Path Transformation function test ------------------------ //

func TestPathTransformFunc(t *testing.T) {
	key := "momo.gif"
	pathKey := CASPathTransformFunc(key)
	fmt.Println(pathKey) // {4295b/3c8f0/a940a/89eb1/0af84/15c9d/8ff32/34234 4295b3c8f0a940a89eb10af8415c9d8ff3234234}
	fmt.Println(pathKey.FullPath())

	expectedFilename := "4295b3c8f0a940a89eb10af8415c9d8ff3234234"

	expectedPathName := "4295b/3c8f0/a940a/89eb1/0af84/15c9d/8ff32/34234"

	if pathKey.Filename != expectedFilename {
		t.Errorf(" Recieved %s expected %s", pathKey.Filename, expectedFilename)
	}
	if pathKey.PathName != expectedPathName {
		t.Errorf("Recieved %s expected %s", pathKey.PathName, expectedPathName)
	}
}

// -------------------------- Store Read, Write, Delete test ------------------------- //
func TestStore(t *testing.T) {
	s := newStore()
	id := "user_momo"
	defer cleanupTest(t, s)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("foo_%d", i)
		data := []byte("some jpg bytes")

		if _, err := s.writeStream(id, key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(id, key); !ok {
			t.Errorf("expected to have key %s", key)
		}

		_, r, err := s.Read(id, key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)
		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}
		fmt.Println("Data: ", string(b))

		if err := s.Delete(id, key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(id, key); ok {
			t.Errorf("expected to NOT have key %s", key)
		}
	}
}

// -------------------- Write Test ------------------------ //

/*
func TestWriteFunc(t *testing.T) {
  // Write the file to the store
  store := newStore()
  defer cleanupTest(t, store)
	id := "user102"
	key := "momo.png"
	data := []byte("some png bytes")

	if _, err := store.Write(id, key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	check := store.Has(id, key)
	assert.True(t, check)
}
*/

// --------------------- Read Test ------------------------- //

func TestReadFunction(t *testing.T) {

	store := newStore()
	id := "user103"
	key := "momo.md"
	data := []byte("markdown data in bytes")
	defer cleanupTest(t, store)
	// write the file with the data
	if _, err := store.Write(id, key, bytes.NewReader(data)); err != nil { // we provide bytes.NewReader to convert the data to io.Reader
		t.Error(err)
	}

	// read the file
	_, reader, err := store.Read(id, key)
	if err != nil {
		t.Error(err)
	}
	b, _ := io.ReadAll(reader)
	fmt.Println("Data: ", string(b))

	// TEST

	assert.Equal(t, data, b)

}

// -------------------- Delete Test ------------------------ //

/*
func TestDeleteFunc(t *testing.T) {
	store := newStore()
	id := "user101"
	key := "momo.jpeg"
  // defer cleanupTest(t, store) // already deleting might hinder assert.False
	data := []byte("some jpg bytes")
	if _, err := s.writeStream(id, key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
  // TEST
	if err := store.Delete(id, key); err != nil {
		t.Error(err)
	}
	assert.False(t, s.Has(id, key)) // this should return false

}
*/
// ----------------------------------- xxxxxxxxxxx ----------------------------------- //
