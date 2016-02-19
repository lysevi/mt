package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/lysevi/mt/storage"
)

var _ = fmt.Sprintf("")

const (
	pingPeriod = time.Duration(500) * time.Millisecond
)

type Server struct {
	is_work    bool
	workers_wg sync.WaitGroup
	port       string
	listen     net.Listener
	Connects   uint32
	clients    []*ClientInfo
	ping_chan  chan interface{}
	Store      *storage.Storage
}

func NewServer(port string) Server {
	s := Server{}
	s.is_work = false
	s.port = port
	s.Connects = 0
	s.ping_chan = make(chan interface{})
	s.Store = storage.NewStorage()
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
	//go s.ping_worker()
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
				break
			}

			panic(fmt.Sprintf("Accept failed: %v", err))
			continue
		} else {
			go s.on_connect(conn)
		}
	}
	s.workers_wg.Done()
	log.Println("server: net_worker stoped")
}

func (s *Server) ping_worker() {
L:
	for {
		time.Sleep(pingPeriod)
		select {
		case <-s.ping_chan:
			break L
		default:
		}

		for i := range s.clients {
			if time.Since(s.clients[i].pingTime) > pingPeriod {
				log.Println("server: ping to ", s.clients[i].conn.LocalAddr().String())
				s.clients[i].conn.Write([]byte(ping))
			}
		}
	}
	log.Println("server: ping_worker stoped")
	s.workers_wg.Done()
}

func (s *Server) client_io_worker(ci *ClientInfo) {
	//log.Println("server: start worker ", ci.conn.LocalAddr().String())
	ci.conn.Write([]byte(helloFromServer))
	buf := make([]byte, 1024, 1024)
	protocol := NewServerProtocol(s)
L:
	for {

		select {
		case <-ci.stop_worker:
			log.Println("server: client_io_worker stoping")
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

		//log.Printf("server: recv: '%v'", strings.Replace(string(buf[:n]), "\n", "<", -1))
		is_close, _ := protocol.OnRecv(ci, buf[:n])
		if is_close {
			s.removeClient(ci)
			break
		}
	}
	s.workers_wg.Done()
	//log.Println("server: stop client_io_worker ", ci.conn.LocalAddr().String())
}

func (s *Server) on_connect(conn net.Conn) {
	//	log.Println("server: on_connect")
	ci := NewClientInfo(conn, s)
	s.clients = append(s.clients, ci)
	s.Connects++
	s.workers_wg.Add(1)
	go s.client_io_worker(ci)
}

func (s *Server) NewQuery(ci *ClientInfo, buf []byte) bool {
	//log.Println("server: new query ", string(buf))
	query_s := string(buf[len(queryRequest):])
	var id int32
	query := make([]byte, 1024, 1024)
	_, err := fmt.Sscanf(query_s, "%d %s", &id, &query)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(ci.conn)
	for i := 0; i < 4; i++ {
		bts, _ := reader.ReadBytes('\n')
		//		log.Println("server: ", string(bts))
		if IsOk(bts) {
			break
		}
		query = append(query, bts...)
	}
	//log.Println("server: new query ", id, string(query))
	for _, v := range s.clients {
		if v.id == id {
			go v.NewQuery(ci, query)
			return true
		}
	}
	log.Panicf("unknow id: %v", id)
	return true
}

func (s *Server) Pong(ci *ClientInfo) bool {
	log.Println("server: pong from ", ci.String())
	ci.pingTime = time.Now()
	return false
}

func (s *Server) SayHello(ci *ClientInfo, buf []byte) bool {
	//	log.Println("server: say hello")
	ci.name = strings.Replace(string(buf), "\n", "", -1)
	log.Printf("server: hello %v", ci.String())
	ci.conn.Write([]byte(fmt.Sprintf("%d", ci.id)))
	return false
}

func (s *Server) Error(ci *ClientInfo, msg string) bool {
	log.Panicln(fmt.Sprint("server: error ", msg))
	return true
}

func (s *Server) Disconnect(ci *ClientInfo) bool {
	log.Println("server: disconnect ", ci.String())
	s.workers_wg.Add(1)
	defer s.workers_wg.Done()

	s.removeClient(ci)

	ci.stoped = true
	close(ci.stop_worker)
	log.Println("server: clients after ", len(s.clients))
	return false
}

func (s *Server) removeClient(ci *ClientInfo) {
	pos := -1
	for i, v := range s.clients {
		if v.id == ci.id {
			pos = i
			break
		}
	}
	if pos == -1 {
		log.Panicf("server: client not found id=%d", ci.String())
	} else {
		s.clients = append(s.clients[:pos], s.clients[pos+1:]...)
	}
}
