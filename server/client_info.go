package server

import (
	"net"
	"sync"
	"time"
)

var client_num = 0

type IOData struct {
	id   int64
	data []byte
}

type ClientInfo struct {
	id          int
	conn        net.Conn
	pingTime    time.Time
	stop_worker chan interface{}
	name        string
	stoped      bool
	mutex       sync.Mutex
	out         []*IOData
}

func NewClientInfo(conn net.Conn) *ClientInfo {
	res := ClientInfo{}
	res.conn = conn
	res.pingTime = time.Now()
	res.stop_worker = make(chan interface{})
	res.id = client_num
	res.stoped = false
	client_num++
	return &res
}

func (c *ClientInfo) addData(d *IOData) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.out = append(c.out, d)
}
