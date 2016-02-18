package server

import (
	"sync"
	"testing"
)

func TestServerStartStop(t *testing.T) {
	serv := NewServer("")
	wg := sync.WaitGroup{}
	wg.Add(1)
	serv.Start()
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
