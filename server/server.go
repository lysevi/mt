package server

import (
	"fmt"
	"sync"
)

var _ = fmt.Sprintf("")

type Server struct {
	stop_chan chan interface{}
	is_work   bool

	workers_wg sync.WaitGroup

	port string
}

func NewServer(port string) Server {
	s := Server{}
	s.is_work = false
	s.stop_chan = make(chan interface{})
	s.port = port
	return s
}

func (s *Server) Start() {
	s.is_work = true
	s.workers_wg.Add(1)
	go s.net_worker()
}

func (s *Server) Stop() {
	s.stop_chan <- 1
}

func (s *Server) net_worker() {
	for {

		select {
		case <-s.stop_chan:
			s.workers_wg.Done()
			s.is_work = false
			break
		}
	}
}
