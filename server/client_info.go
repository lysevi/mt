package server

import (
	"net"
	"time"
)

var client_num = 0

type ClientInfo struct {
	id          int
	conn        net.Conn
	pingTime    time.Time
	stop_worker chan interface{}
	name        string
	stoped      bool
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
