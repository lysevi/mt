package main

import (
	"log"

	_ "github.com/lysevi/mt/storage"
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
}
