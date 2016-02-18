package server

import (
	"fmt"
	"net"
	"sync"
)

var _ = fmt.Sprintf("")

type Server struct {
	is_work bool

	workers_wg sync.WaitGroup

	port   string
	listen net.Listener

	Connects uint32
}

func NewServer(port string) Server {
	s := Server{}
	s.is_work = false
	s.port = port
	s.Connects = 0
	return s
}

func (s *Server) Start() error {
	var err error
	s.is_work = true
	s.workers_wg.Add(1)
	s.listen, err = net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	go s.net_worker()
	return nil
}

func (s *Server) Stop() {
	s.listen.Close()
}

func (s *Server) Wait() {
	s.workers_wg.Wait()
}

func (s *Server) net_worker() {
	for {
		conn, err := s.listen.Accept()
		if err != nil {
			if x, ok := err.(*net.OpError); ok && x.Op == "accept" { // We're done
				s.is_work = false
				s.workers_wg.Done()
				break
			}

			panic(fmt.Sprintf("Accept failed: %v", err))
			continue
		} else {
			go s.on_connect(conn)
		}
	}
}

func (s *Server) on_connect(conn net.Conn) {
	fmt.Println("server: on_connect")
	s.Connects++
	conn.Write([]byte("+++\n"))
	conn.Close()
}
