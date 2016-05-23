package gokvstore

import (
	"testing"
	"sync"
	"fmt"
)

func TestStore_Init(t *testing.T) {
	var s *Store
	s = &Store{}
	s.Init()
	if s.storage == nil {
		t.Error("Expected s.storage to not be nil")
	}
}

func TestStore_WriteItem(t *testing.T) {
	s := &Store{}
	s.Init()

	s.WriteItem(storageItem{Key: "test", Value: true})
	if val, _ := s.GetItem(storageItem{Key: "test"}); val != true {
		t.Error("Expected s.GetItem(test) to return true, got", val)
	}
}

func TestStore_GetItem(t *testing.T) {
	s := &Store{}
	s.Init()

	if _, err := s.GetItem(storageItem{Key: "wrong"}); err == nil {
		t.Error("Expected s.GetItem(wrong) to throw an error with wrong key, got", err)
	}

	s.WriteItem(storageItem{Key: "test", Value: true})
	if _, err := s.GetItem(storageItem{Key: "test"}); err != nil {
		t.Error("Expected s.GetItem(test) to throw NO error with wrong key, got", err)
	}
}

func TestConcurrent(t *testing.T) {
	s := &Store{}
	s.Init()

	wg := &sync.WaitGroup{}

	wg.Add(3)

	// Run 3 writes in goroutines, which should run in order thanks to the mutex
	go func() {
		s.WriteItem(storageItem{"test", 1})
		wg.Done()
	}()
	go func() {
		s.WriteItem(storageItem{"test", 2})
		wg.Done()
	}()
	go func() {
		s.WriteItem(storageItem{"test", 3})
		wg.Done()
	}()

	wg.Wait()
	if val, _ := s.GetItem(storageItem{Key: "test"}); val != 3 {
		t.Error("val != 3, got", val)
	}
}