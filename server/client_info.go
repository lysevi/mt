package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var client_num int32 = 1

type IOData struct {
	id   int64
	data []byte
}

type ClientInfo struct {
	id          int32
	conn        net.Conn
	pingTime    time.Time
	stop_worker chan interface{}
	name        string
	stoped      bool
	mutex       sync.Mutex
	out         []*IOData

	queryes int
}

func NewClientInfo(conn net.Conn) *ClientInfo {
	res := ClientInfo{}
	res.conn = conn
	res.pingTime = time.Now()
	res.stop_worker = make(chan interface{})
	res.id = client_num
	res.name = "error name"
	res.stoped = false
	client_num++
	return &res
}

func (c *ClientInfo) String() string {
	return fmt.Sprintf("{id:%v name:'%v'}", c.id, c.name)
}

func (c *ClientInfo) addData(d *IOData) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.out = append(c.out, d)
}

func (c *ClientInfo) NewQuery(queryClient *ClientInfo, buf []byte) {
	log.Println("server: new query ", c.String(), "Q=", string(buf[:len(buf)-1]))
	c.queryes++
	queryClient.conn.Write([]byte(ok))
}
