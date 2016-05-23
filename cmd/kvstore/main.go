package main

import (
	"github.com/wizardfrag/gokvstore"
	"sync"
)

func main() {
	s := &gokvstore.Store{}
	s.Init()

	// Use a channel so we automatically exit if we can't bind to one of the ports
	shutdownChan := make(chan bool)

	// Start TCP Server
	go func() {
		(&gokvstore.TcpServer{}).Init(6000, s)
		shutdownChan <- true
	}()
	// Start UDP Server
	go func() {
		(&gokvstore.UdpServer{}).Init(6000, s)
		shutdownChan <- true
	}()

	// Start HTTP Server
	go func() {
		(&gokvstore.HttpServer{}).Init(8000, s)
		shutdownChan <- true
	}()

	<- shutdownChan
}
