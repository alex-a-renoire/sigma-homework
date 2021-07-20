package server

import (
	"log"
	"net"
	"sync"

	"github.com/alex-a-renoire/tcp/tcpserver/handler"
)

type TCPServer struct {
	Listener net.Listener
	Handler  handler.Handler
	Quit     chan interface{}
	Wg       sync.WaitGroup
}

func New(addr string, handler handler.Handler) *TCPServer {
	// creating a TCP listener
	log.Print("Server starting...")

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error listening on %s: %s", addr, err)
	}

	return &TCPServer{
		Quit:     make(chan interface{}),
		Listener: listener,
		Handler:  handler,
	}
}

func (s *TCPServer) Serve() {
	defer s.Wg.Done()

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			select {
			case <-s.Quit:
				return
			default:
				log.Fatalf("error establishing connection: %s", err)
			}
		}
		log.Print("connection accepted")

		s.Wg.Add(1)
		go func() {
			s.Handler.HandleConnection(conn)
			s.Wg.Done()
		}()
	}
}

func (s *TCPServer) Stop() {
	close(s.Quit)
	s.Listener.Close()
	s.Wg.Wait()
}
