//TCP server lifecycle + connection handling

package server

import (
	"bufio"
	"fmt"
	"net"
)

// Server is the TCP entrypoint
type Server struct {
	addr    string
	handler *Handler
}

// NewServer creates server instance
func NewServer(addr string, handler *Handler) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
	}
}

// Start launches TCP server
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Println("Horcrux TCP server running on", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go s.handleConnection(conn) // concurrency model
	}
}

// handleConnection processes one client
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	reader := bufio.NewScanner(conn)

	for reader.Scan() {
		line := reader.Text()

		fmt.Println("Received command:", line)

		cmd := ParseCommand(line)
		resp := s.handler.Execute(cmd)

		fmt.Println("Sending response:", resp)

		fmt.Fprintln(conn, resp)
	}

	fmt.Println("Client disconnected:", conn.RemoteAddr())
}
