package gokvstore

import (
	"sync"
	"errors"
)

// Store implements the key-value store
type Store struct {
	storage kvStore
	storageMutex sync.Mutex
}

func (s *Store) Init() {
	s.storage = make(kvStore)
}

func (s *Store) WriteItem(item storageItem) {
	s.storageMutex.Lock()
	defer s.storageMutex.Unlock()
	s.storage[item.Key] = item.Value
}

func (s *Store) GetItem(item storageItem) (kvItem, error) {
	s.storageMutex.Lock()
	defer s.storageMutex.Unlock()
	if val, ok := s.storage[item.Key]; ok {
		return val, nil
	}
	return nil, errors.New("No such key")
}