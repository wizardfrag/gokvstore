package main

import (
	"sync"
	"github.com/wizardfrag/gokvstore"
)

func main() {
	s := &gokvstore.Store{}
	s.Init()
	var wg sync.WaitGroup

	// Use a sync.WaitGroup to wait for all servers to finish...
	wg.Add(3)

	// Start TCP Server
	go func() {
		defer wg.Done()
		(&gokvstore.TcpServer{}).Init(6000, s)
	}()
	// Start UDP Server
	go func() {
		defer wg.Done()
		(&gokvstore.UdpServer{}).Init(6000, s)
	}()

	// Start HTTP Server
	go func() {
		defer wg.Done()
		(&gokvstore.HttpServer{}).Init(8000, s)
	}()

	wg.Wait()
}