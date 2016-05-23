package gokvstore

import (
	"fmt"
	"net"
	"encoding/json"
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

	s.conn, err = net.ListenUDP("udp", addr)

	if err != nil {
		panic(err)
	}
	defer s.conn.Close()
	fmt.Println("UDP Server listening on port", port)
	for {
		var buf = make([]byte, 1024)
		n, address, err := s.conn.ReadFromUDP(buf)

		s.handleRequest(address, string(buf[0:n]))

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

func (s *UdpServer) handleRequest(address *net.UDPAddr, buf string) {
	var command storageCommand
	err := json.Unmarshal([]byte(buf), &command)
	if err != nil {
		s.writeResult(address, &kvError{err.Error(), 1})
		return
	}
	if command.Cmd == "put" {
		s.store.WriteItem(command.Item)
		s.writeResult(address, &kvResponse{Success: true})
	} else if command.Cmd == "get" {
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

func (s *UdpServer) writeResult(address *net.UDPAddr, data interface{}) {
	jsonBuf, _ := json.Marshal(data)
	jsonBuf = append(jsonBuf, '\n')
	s.conn.WriteToUDP(jsonBuf, address)
}