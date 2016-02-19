package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lysevi/mt/server"
	"github.com/lysevi/mt/storage"
	//"github.com/pkg/profile"
)

type emptyLogger int

func (c emptyLogger) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var el emptyLogger

func init() {
	//log.SetOutput(el)
}

func main() {
	//defer profile.Start().Stop()
	log.Println("****************")

	serv := server.NewServer(":8080")
	serv.Start()
	wg := sync.WaitGroup{}
	wg.Add(1)
	start_time := time.Now()
	var elapsed_time time.Duration
	f := func(name string, start *time.Time, elapsed *time.Duration) {
		conn, err := server.Connect(name, "localhost:8080")
		if err != nil {
			panic(err)
		}

		vals := []server.Value{}
		for i := 0; i < 100000; i++ {
			v := server.Value{Id: 0, Time: storage.Time(i), Value: int64(i), Flag: 0xff}
			vals = append(vals, v)
		}
		*start = time.Now()
		conn.WriteValues(vals)
		*elapsed = time.Now().Sub(*start)

		conn.Disconnect()
		wg.Done()
	}

	go f("client 1", &start_time, &elapsed_time)

	wg.Wait()
	serv.Stop()

	fmt.Println("************ elapsed: ", elapsed_time)
}
