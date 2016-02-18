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

	go func() {
		_, err := server.Connect("main.go", "localhost:8080")
		if err != nil {
			panic(err)
		}

		for {
		}
	}()

	for {
	}
}
