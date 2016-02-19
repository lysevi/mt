package server

import (
	"log"

	"github.com/lysevi/mt/storage"
)

type Writer struct {
	stor *storage.Storage
	stop chan interface{}
}

func NewWriter(stor *storage.Storage) *Writer {
	res := Writer{}
	res.stor = stor
	res.stop = make(chan interface{})
	return &res
}

func (w *Writer) addValue(v Value) {
	w.stor.Add(storage.NewMeas(v.Id, v.Time, v.Value, v.Flag))
}

func (w *Writer) Stop() {
	close(w.stop)
	log.Println("server: writer stoped")
}
