package server

import (
	"sync"
	"testing"
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

	client, err := Connect("localhost:8080")
	if err != nil {
		t.Error("client connect error")
		return
	}

	for {
		if serv.Connects == 1 && client.is_connected {
			break
		}
	}

	client.Disconnect()

	if client.is_connected || !client.is_closed {
		t.Error("client close error: ", client.is_connected, client.is_closed)
	}
	serv.Stop()

}
