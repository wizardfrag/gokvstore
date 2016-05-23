package gokvstore

import (
	"fmt"
	"net"
	"encoding/json"
	"strings"
)

type UdpServer struct {
	conn *net.UDPConn
	store *Store
}

func (s *UdpServer) Init(port int, st *Store) {
	var err error
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	s.store = st

	if err != nil {
		panic(err)
	}

	// create the listener
	s.conn, err = net.ListenUDP("udp", addr)

	if err != nil {
		panic(err)
	}

	defer s.conn.Close()
	fmt.Println("UDP Server listening on port", port)

	for {
		// loop, listening for inbound packets
		// This assumes that one packet contains one command, and clients should
		// be written to assume this also
		var buf = make([]byte, 1024)
		n, address, err := s.conn.ReadFromUDP(buf)


		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		s.handleRequest(address, string(buf[0:n]))
	}
}

func (s *UdpServer) handleRequest(address *net.UDPAddr, buf string) {
	var command storageCommand
	err := json.Unmarshal([]byte(buf), &command)
	if err != nil {
		s.writeResult(address, &kvError{err.Error(), 1})
		return
	}

	// Check the command we've been given
	if strings.ToLower(command.Cmd) == "put" {
		s.store.WriteItem(command.Item)
		// Send back a success: true, this will always be true, pretty much
		s.writeResult(address, &kvResponse{Success: true})
	} else if strings.ToLower(command.Cmd) == "get" {
		item, err := s.store.GetItem(command.Item)
		if err != nil {
			s.writeResult(address, &kvError{err.Error(), 2})
			return
		}
		s.writeResult(address, item)
	} else {
		s.writeResult(address, &kvError{"Invalid command", 3})
	}
}

// writeResult is a simple implementation similar to HttpServer.httpOutput
func (s *UdpServer) writeResult(address *net.UDPAddr, data interface{}) {
	jsonBuf, _ := json.Marshal(data)
	jsonBuf = append(jsonBuf, '\n')
	s.conn.WriteToUDP(jsonBuf, address)
}