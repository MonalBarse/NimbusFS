package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "nimbud_root"

// --------------------------Path Transform Functions------------------------------------ //
func CASPathTransformFunc(key string) PathKey {
	// Hash the key using SHA-1
	hash := sha1.Sum([]byte(key))          // sha1 gives out a hash of 20 bytes of the key. Lets say the key is "How are we doing?" then the hash will be "6804429f74181a63c50c3d81d733a12f14a353ff" of length 40
	hashStr := hex.EncodeToString(hash[:]) // we had to do hex encoding to convert the hash to a string as the hash is a byte array so first the hex encoding is done and then the string conversion is done

	// Let's say the hash is 6804429f74181a63c50c3d81d733a12f14a353ff so with the blocksize of 5 we will get the following paths : 68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff
	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"), // Concatenate hash parts to form a path
		Filename: hashStr,                  // The full hash as the filename
	}
}

type PathTransformFunc func(string) PathKey

// Default path is provided because if user does not provide path and filename and just provides the content 'key' then we can use this default path to store the content
var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key, // Default path transformation just uses the key
		Filename: key,
	}
}

// ------------------------------ XXXXXXXXXXXXXXX---------------------------------------- //
// ------------------------------- PathKey Struct --------------------------------------- //
type PathKey struct {
	PathName string // The path formed from the hashed key
	Filename string // The full hash used as the filename
} // PathKey might look like if user provided key as "How you doin'?" the filename will be "6804429f74181a63c50c3d81d733a12f14a353ff" and the path will be "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff

func (p PathKey) FirstPathName() string {
	// Retrieve the first part of the path
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath() string {
	// Construct the full path by combining PathName and Filename
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
} // FullPath then will be "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff/6804429f74181a63c50c3d81d733a12f14a353ff"

// ------------------------------ XXXXXXXXXXXXXXX---------------------------------------- //
//------------------------ Store Options and Constructor -------------------------------- //

type Store struct {
	StoreOpts // Embedding StoreOpts to use its fields directly
}

type StoreOpts struct {
	Root              string            // Root folder for the storage system
	PathTransformFunc PathTransformFunc // Function to transform keys into paths
}

func NewStore(opts StoreOpts) *Store {
	// Initialize default path transform function if not provided
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	// Set default root if not provided
	if opts.Root == "" {
		opts.Root = defaultRootFolderName
	}

	return &Store{
		StoreOpts: opts,
	}
}

// Let's say the root is "nimus_dir" and the path transform function is CASPathTransformFunc then the store will be created with these options
// and for eg. if the user provides : key = 'hello.txt' id = 'user1' then the path will be: nimus_dir/user1/68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff/hello.txt

// ------------------------------ XXXXXXXXXXXXXXX---------------------------------------- //
// --------------------------------Store Methods ---------------------------------------- //

/* 1. Write the file to the store and return the number of bytes written ---------------- */
func (s *Store) Write(id string, key string, r io.Reader) (int64, error) {
	// Write data to a file in the store
	return s.writeStream(id, key, r)
}

func (s *Store) writeStream(id string, key string, r io.Reader) (int64, error) {
	// Write data from the reader to the file
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	return io.Copy(f, r)
}

func (s *Store) openFileForWriting(id string, key string) (*os.File, error) {
	// Open a file for writing, creating directories as needed
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.PathName) // The pathNameWithRoot will be "nimus_dir/user1/68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {        // os.MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
		return nil, err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath()) // The fullPathWithRoot will be "nimus_dir/user1/68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff/hello.txt"

	return os.Create(fullPathWithRoot)
}

/* 2. Read the file from the store and return the number of bytes read and the reader --- */

func (s *Store) Read(id string, key string) (int64, io.Reader, error) {
	// Read data from a file in the store
	return s.readStream(id, key)
}

func (s *Store) readStream(id string, key string) (int64, io.ReadCloser, error) {
	// Read data from the file and return the size and reader
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	file, err := os.Open(fullPathWithRoot)
	if err != nil {
		return 0, nil, err
	}

	fileStat, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}
	sizeOfFile := fileStat.Size()

	return sizeOfFile, file, nil
}

/* 3. Check if the file exists in the store ---------------------------------------------- */
func (s *Store) Has(id string, key string) bool {
	// Check if a file exists in the store
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	_, err := os.Stat(fullPathWithRoot)    // os.Stat returns file info. It will return an error if the file does not exist
	return !errors.Is(err, os.ErrNotExist) // errors.Is reports whether any error in err's chain matches target. os.ErrNotExist is the error returned by os.Stat when the file does not exist
}

/* 4. Delete the file from the store ----------------------------------------------------- */

func (s *Store) Delete(id string, key string) error {
	// Delete a specific file from the store
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()

	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FirstPathName())

	return os.RemoveAll(pathNameWithRoot) // RemoveAll removes path and any children it contains.
}

// for testing purposes to clear away the entire storage
func (s *Store) Clear() error {
	// Clear the entire storage by removing the root directory
	return os.RemoveAll(s.Root)
}

// ------------------------------ XXXXXXXXXXXXXXX----------------------------------------- //
