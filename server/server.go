package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var _ = fmt.Sprintf("")

const (
	pingPeriod = time.Duration(5) * time.Second
)

type Server struct {
	is_work    bool
	workers_wg sync.WaitGroup
	port       string
	listen     net.Listener
	Connects   uint32
	clients    []*ClientInfo

	ping_chan chan interface{}
}

func NewServer(port string) Server {
	s := Server{}
	s.is_work = false
	s.port = port
	s.Connects = 0
	s.ping_chan = make(chan interface{})
	return s
}

func (s *Server) Start() error {
	var err error
	s.is_work = true

	s.listen, err = net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	s.workers_wg.Add(1)
	go s.net_worker()
	s.workers_wg.Add(1)
	go s.ping_worker()
	return nil
}

func (s *Server) Stop() {
	log.Println("server: stoping ")
	close(s.ping_chan)
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

func (s *Server) ping_worker() {
L:
	for {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		select {
		case <-s.ping_chan:
			break L
		default:
		}
		//log.Println("server: pings")
		for i := range s.clients {
			if time.Since(s.clients[i].pingTime) > pingPeriod {
				log.Println("server: ping to ", s.clients[i].conn.LocalAddr().String())
				s.clients[i].pingTime = time.Now()
			}
		}
	}
	s.workers_wg.Done()
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

		log.Printf("server: recv: '%v'", strings.Replace(string(buf[:n]), "\n", "<", -1))
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
	ci.name = strings.Replace(string(buf), "\n", "<", -1)
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
