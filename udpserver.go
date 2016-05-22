package gokvstore

import (
	"net"
	"fmt"
)

type UdpServer struct {
	conn *net.UDPConn
}

func (s *UdpServer) Init(port int, st *Store) {
	var err error
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
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
		fmt.Println("Received ", string(buf[0:n]), " from ", address)

		if err != nil {
			fmt.Println("Error: ",err)
		}
	}
}
