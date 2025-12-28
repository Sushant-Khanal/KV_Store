package main

import (
	"fmt"
	"log"
	"sync"
)

type Storer[K comparable, V any] interface {
	Put(K, V) error
	Get(K) (V, error)
	Update(K, V) error
	Delete(K) (V, error)
}

type KVStore[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	return &KVStore[K, V]{
		data: make(map[K]V),
	}
}

// Has checks if the given key is present in the store.
// Note: This is not Concurrent Safe, should be used with lock or mutex.
func (s *KVStore[K, V]) Has(key K) bool {
	_, ok := s.data[key]
	return ok
}

func (s *KVStore[K, V]) Update(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Has(key) {
		return fmt.Errorf("The key (%v) does not exists", value)
	}

	s.data[key] = value

	return nil
}

func (s *KVStore[K, V]) Put(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return nil
}

func (s *KVStore[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return value, fmt.Errorf("the key (%v) does not exists", key)
	}

	return value, nil
}

type Server struct {
	Store Storer[string, string]
}

func (s *Server) getUserByName(name string) (string, error) {
	return s.Store.Get(name)

}

// type Block struct{}
// type Transaction struct{}

func main() {
	store := NewKVStore[string, string]()

	if err := store.Put("foo", "bar"); err != nil {
		log.Fatal(err)
	}

	value, err := store.Get("foo")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(value)

	if err := store.Put("foo", "oof"); err != nil {
		log.Fatal(err)
	}
	value, err = store.Get("foo")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(value)

	//StoreThings(kv)
}
