package server

import (
	"encoding/json"
	_ "log"

	"github.com/lysevi/mt/storage"
)

const (
	queryWriteKind = "write"
	queryReadKind  = "read"

	queryWrite = 1 << iota
	queryRead  = iota
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
	res.Kind = queryWriteKind
	return res
}

func NewQueryRead() *QueryRead {
	res := &QueryRead{}
	res.Kind = queryReadKind
	return res
}

func (q QueryWrite) bytes() ([]byte, error) {
	res, err := json.Marshal(q)
	//log.Println("marshal: ", q, string(res))
	return res, err
}

func (q QueryRead) bytes() ([]byte, error) {
	res, err := json.Marshal(q)
	return res, err
}
