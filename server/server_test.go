package server

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestServerStartStop(t *testing.T) {
	serv := NewServer("")
	wg := sync.WaitGroup{}
	wg.Add(1)

	if err := serv.Start(); err != nil {
		t.Error("start error: ", err)
		return
	}
	go func(s *Server, w *sync.WaitGroup) {
		for {
			if !serv.is_work {
				break
			}
		}
		wg.Done()
	}(&serv, &wg)

	serv.Stop()
	wg.Wait()
}

func TestServerConnect(t *testing.T) {
	serv := NewServer(":8080")

	if err := serv.Start(); err != nil {
		t.Error("start error: ", err)
		return
	}

	client, err := Connect("test", "localhost:8080")
	if err != nil {
		t.Error("client connect error")
		return
	}

	for {
		if serv.Connects == 1 && client.is_connected {
			break
		}
	}
	time.Sleep(time.Duration(500) * time.Millisecond)

	client.Disconnect()
	if client.is_connected || !client.is_closed {
		t.Error("client close error: ", client.is_connected, client.is_closed)
	}

	serv.Stop()
}

func TestServerClientQuerys(t *testing.T) {
	serv := NewServer(":8080")
	serv.Start()

	f := func(name string) {
		conn, err := Connect(name, "localhost:8080")
		if err != nil {
			panic(err)
		}
		conn.SendQuery([]byte("test query 1"))
		conn.SendQuery([]byte("test query 2"))
		log.Println("client: ", name, " end")
		conn.Disconnect()
	}

	go f("client 1")
	go f("client 2")

	time.Sleep(time.Duration(2) * time.Second)
	serv.Stop()
}
