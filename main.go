package main

import (
	"fmt"

	"github.com/lysevi/mt/server"
	//	"github.com/lysevi/mt/storage"
	//"github.com/pkg/profile"
)

func main() {
	//defer profile.Start().Stop()
	fmt.Println("****************")

	serv := server.NewServer(":8080")
	serv.Start()

	f := func(name string) {
		conn, err := server.Connect(name, "localhost:8080")
		if err != nil {
			panic(err)
		}
		conn.SendQuery([]byte("test query 1"))
		conn.SendQuery([]byte("test query 2"))
	}

	go f("client 1")
	go f("client 2")

	for {
	}
}
