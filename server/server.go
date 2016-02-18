package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var _ = fmt.Sprintf("")

type Server struct {
	is_work    bool
	workers_wg sync.WaitGroup
	port       string
	listen     net.Listener
	Connects   uint32
	clients    []*ClientInfo
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
	log.Println("server: stoping ")
	s.listen.Close()
	for i := range s.clients {
		log.Println("server: close ", s.clients[i].name)
		s.clients[i].stop_worker <- 1
	}
	log.Println("server: stop wait")
	s.Wait()
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

func (s *Server) client_io_worker(ci *ClientInfo) {
	log.Println("server: start worker ", ci.conn.LocalAddr().String())
	ci.conn.Write([]byte(helloFromServer))
	buf := make([]byte, 1024, 1024)
	protocol := NewServerProtocol(s)
L:
	for {

		select {
		case <-ci.stop_worker:
			log.Println("server: worker stoping")
			break L
		default:
		}
		ci.conn.SetDeadline(time.Now().Add(time.Duration(500) * time.Millisecond))
		n, err := ci.conn.Read(buf)
		if err != nil {
			opErr, ok := err.(*net.OpError)
			if ok && (opErr.Timeout() || opErr.Err == io.EOF) {
				continue
			}
			if err != io.EOF {
				log.Println("worker error: ", err)
				break L
			}
		}
		if n == 0 {
			continue
		}

		log.Printf("server: recv: '%v'", string(buf[:n]))
		protocol.OnRecv(ci, buf[:n])
	}
	s.workers_wg.Done()
	log.Println("server: stop worker ", ci.conn.LocalAddr().String())
}

func (s *Server) on_connect(conn net.Conn) {
	log.Println("server: on_connect")
	ci := NewClientInfo(conn)
	s.clients = append(s.clients, ci)
	s.Connects++
	s.workers_wg.Add(1)
	go s.client_io_worker(ci)
}

func (s *Server) Pong(ci *ClientInfo) {
	log.Println("server: pong from ", ci.conn.LocalAddr().String())
	ci.pingTime = time.Now()
}

func (s *Server) SayHello(ci *ClientInfo, buf []byte) {
	log.Println("server: say hello")
	ci.name = string(buf)
	log.Printf("server: hello %v, id=%d", ci.name, ci.id)

}

func (s *Server) Error(ci *ClientInfo, msg string) {
	log.Panicln(fmt.Sprint("server: error ", msg))
}

func (s *Server) Disconnect(ci *ClientInfo) {
	log.Println("server: disconnect ", ci.conn.LocalAddr())
	s.workers_wg.Add(1)
	defer s.workers_wg.Done()
	pos := -1
	for i, v := range s.clients {
		if v.id == ci.id {
			pos = i
			break
		}
	}
	if pos == -1 {
		log.Panicf("server: client not found id=%d", ci.id)
	}

	s.clients = append(s.clients[:pos], s.clients[pos+1:]...)
	ci.stoped = true
	close(ci.stop_worker)
	log.Println("server: clients ", len(s.clients))

}
