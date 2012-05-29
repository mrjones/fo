package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Function which can fetch a resource.  Returns the contents of that resource
// (as a string) if successful, and an error otherwise.
type FetchFunction func() (body string, err error)

type ReadThroughCache struct {
	storage KVStore
}

func NewReadThroughCache(storage KVStore) ReadThroughCache {
	return ReadThroughCache{storage: storage}
}

func (c ReadThroughCache) GetAsReader(
	  ff FetchFunction, cachekey string, maxage time.Duration) (io.Reader, error) {
	contents, err := c.Get(ff, cachekey, maxage)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(contents), nil
}

func (c ReadThroughCache) Get(ff FetchFunction, cachekey string, maxage time.Duration) (string, error) {
	age := c.storage.Age(cachekey)
	if age == nil || *age > maxage {
		response, err := ff()
		if err != nil {
			return "", err
		}
		err = c.storage.Put(cachekey, response)
		return response, err
	}

	return c.storage.Get(cachekey)
}

type KVStore interface {
	Put(k,v string) error
	Get(k string) (v string, err error)
	Age(k string) *time.Duration
}

//
// FileKVStore
//

type FileKVStore struct {
	rootDir string
}

func NewFileKVStore(rootDir string) FileKVStore {
	return FileKVStore{rootDir: rootDir}
}

func (s FileKVStore) filename(k string) string {
	return s.rootDir + "/" + k
}

func (s FileKVStore) Put(k, v string) error {
	return ioutil.WriteFile(s.filename(k), []byte(v), 0644)
}

func (s FileKVStore) Get(k string) (v string, err error) {
	bytes, err := ioutil.ReadFile(s.filename(k))
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (s FileKVStore) Age(k string) *time.Duration {
	return fileAge(s.filename(k))
}

func fileAge(filename string) *time.Duration {
 	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return nil
	}

	age := time.Since(fileinfo.ModTime())

	return &age
}

//
// MemKVStore
//

type MemKVStore struct {
	kvs map[string]string
	ages map[string]time.Time
}

func NewMemKVStore() MemKVStore {
	return MemKVStore{
		kvs: make(map[string]string),
		ages: make(map[string]time.Time),
	}
}

func (c MemKVStore) Put(k, v string) error {
	c.kvs[k] = v
	c.ages[k] = time.Now()

	return nil
}

func (c MemKVStore) Get(k string) (v string, err error) {
	return c.kvs[k], nil
}

func (c MemKVStore) Age(k string) *time.Duration {
	creation, exists := c.ages[k]
	if exists {
		age := time.Since(creation)
		return &age
	}
	return nil
}

//
// Handy utilities
//

func httpGetBody(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		return "", err
	}

	bits, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(bits), nil
}
