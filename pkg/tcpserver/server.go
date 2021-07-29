package tcpserver

import (
	"log"
	"net"
	"sync"

	"github.com/alex-a-renoire/sigma-homework/pkg/tcpserver/handler"
)

type TCPServer struct {
	Listener net.Listener
	Handler  handler.Handler
	Quit     chan interface{}
	Wg       sync.WaitGroup
	connCounter int
}

func New(addr string, handler handler.Handler) *TCPServer {
	// creating a TCP listener
	log.Print("Server starting...")

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error listening on %s: %s", addr, err)
	}
	log.Print("Server started...")

	return &TCPServer{
		Quit:     make(chan interface{}),
		Listener: listener,
		Handler:  handler,
		connCounter: 0,
	}
}

func (s *TCPServer) Serve() {
	defer s.Wg.Done()

	for {
		s.connCounter++
		numClient := s.connCounter
		conn, err := s.Listener.Accept()
		if err != nil {
			select {
			case <-s.Quit:
				log.Print("Server terminated")
				return
			default:
				log.Fatalf("error establishing connection: %s", err)
			}
		}
		log.Print("connection accepted")

		s.Wg.Add(1)

		ch := make (chan string)

		go s.Handler.HandleConnection(conn, ch, numClient)
		go func() {
			s.Handler.WriterToServer(conn, ch, s.Quit, numClient)
			s.Wg.Done()
		}()
	}
}

func (s *TCPServer) Stop() {
	close(s.Quit)
	s.Listener.Close()
	s.Wg.Wait()
}
