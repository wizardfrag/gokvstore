package gokvstore

import (
	"fmt"
	"log"
	"net"
)

type TcpServer struct {
	listener net.Listener
}

func (s *TcpServer) Init(port int, st *Store) {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	defer s.listener.Close()
	fmt.Println("TCP Server listening on port", port)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Accept error: ", err.Error())
			continue
		}
		go s.handleRequest(conn, st)
	}
}

func (s *TcpServer) handleRequest(conn net.Conn, st *Store) {
	buf := make([]byte, 1024)
	len, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Read error: ", err.Error())
	}
	fmt.Printf("Read %d bytes: %s\n", len, string(buf))
}
