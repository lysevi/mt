package server

import (
	"encoding/json"
	_ "log"

	"github.com/lysevi/mt/storage"
)

const (
	queryWrite = "write"
	queryRead  = "read"
)

type Value struct {
	Id    storage.Id
	Time  storage.Time
	Value int64
	Flag  storage.Flag
}

type QueryRead struct {
	Kind string
	From storage.Time
	To   storage.Time
}

type QueryWrite struct {
	Kind   string
	Values []Value
}

func NewQueryWrite() *QueryWrite {
	res := &QueryWrite{}
	res.Kind = queryWrite
	return res
}

func NewQueryRead() *QueryRead {
	res := &QueryRead{}
	res.Kind = queryRead
	return res
}

func (q QueryWrite) JSON() ([]byte, error) {
	res, err := json.Marshal(q)
	//log.Println("marshal: ", q, string(res))
	return res, err
}

func (q QueryRead) JSON() ([]byte, error) {
	res, err := json.Marshal(q)
	return res, err
}
