package gokvstore

import (
	"fmt"
	"log"
	"net"
	"bufio"
	"os"
	"encoding/json"
	"strings"
)

type TcpServer struct {
	listener net.Listener
	store *Store
}

func (s *TcpServer) Init(port int, st *Store) {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	s.store = st

	if err != nil {
		panic(err)
	}
	defer s.listener.Close()

	fmt.Println("TCP Server listening on port", port)
	// Loop, accepting tcp requests
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Accept error: ", err.Error())
			continue
		}
		go s.handleRequest(conn)
	}
}

func (s *TcpServer) handleRequest(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	// Use scanner.Scan() to retrieve a line at a time
	for scanner.Scan() {
		var command storageCommand
		str := scanner.Text()

		err := json.Unmarshal([]byte(str), &command)
		if err != nil {
			s.writeResult(conn, &kvError{err.Error(), 1})
			continue
		}

		// Check whether this is a get or a write/put command
		if strings.ToLower(command.Cmd) == "put" {
			s.store.WriteItem(command.Item)
			s.writeResult(conn, &kvResponse{Success: true})
		} else if strings.ToLower(command.Cmd) == "get" {
			item, err := s.store.GetItem(command.Item)
			if err != nil {
				s.writeResult(conn, &kvError{err.Error(), 2})
				return
			}

			s.writeResult(conn, item)
		} else {
			s.writeResult(conn, &kvError{"Invalid command", 3})
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading ", conn, "-", err)
	}
}

// writeResult is a simple implementation similar to HttpServer.httpOutput
func (s *TcpServer) writeResult(conn net.Conn, data interface{}) {
	jsonBuf, _ := json.Marshal(data)
	jsonBuf = append(jsonBuf, '\n')
	conn.Write(jsonBuf)
}