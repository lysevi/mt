package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/lysevi/mt/storage"
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
	//	mutex       sync.Mutex
	queryes int
	serv    *Server
}

func NewClientInfo(conn net.Conn, serv *Server) *ClientInfo {
	res := ClientInfo{}
	res.conn = conn
	res.pingTime = time.Now()
	res.stop_worker = make(chan interface{})
	res.id = client_num
	res.name = "error name"
	res.stoped = false
	res.serv = serv
	client_num++
	return &res
}

func (c *ClientInfo) String() string {
	return fmt.Sprintf("{id:%v name:'%v'}", c.id, c.name)
}

func (c *ClientInfo) NewQuery(queryClient *ClientInfo, buf []byte) {
	c.queryes++
	log.Println("server: new query ", c.String(), "Q=", string(buf))

	//TODO rewrite
	qwrite := QueryWrite{}
	err := json.Unmarshal(buf, &qwrite)
	if err == nil && qwrite.Kind == queryWrite {
		log.Println("server: write ", qwrite.Values)
		for _, v := range qwrite.Values {
			c.serv.Store.Add(storage.NewMeas(v.Id, v.Time, v.Value, v.Flag))
		}
		queryClient.conn.Write([]byte(ok))
	}

	qread := QueryRead{}
	err = json.Unmarshal(buf, &qread)
	if err == nil && qread.Kind == queryRead {
		log.Println("server: read ", qread)
		read_res := c.serv.Store.Read([]storage.Id{}, qread.From, qread.To)

		answer := []Value{}
		for _, v := range read_res {
			answer = append(answer, Value{Id: v.Id, Value: v.Value, Time: v.Tstamp})
		}

		answer_json, err := json.Marshal(answer)
		if err == nil {
			answer_str := fmt.Sprintf("%s\n", string(answer_json))
			queryClient.conn.Write([]byte(answer_str))
			queryClient.conn.Write([]byte(ok))
		}
	}

	if err != nil {
		panic(err)
	}
}
