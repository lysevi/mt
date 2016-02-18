package server

import (
	"bufio"
	"fmt"
	"net"
)

var _ = fmt.Sprintf("")

type Client struct {
	conn      net.Conn
	is_closed bool
}

func Connect(con_str string) (*Client, error) {
	conn, err := net.Dial("tcp", con_str)
	if err != nil {
		return nil, err
	}
	res := Client{}
	res.conn = conn
	res.is_closed = false
	go res.client_worker()
	return &res, nil
}

func (c *Client) Disconnect() {
	c.is_closed = true
	c.conn.Close()
}
func (c *Client) client_worker() {
	reader := bufio.NewReader(c.conn)
	for {
		res, err := reader.ReadString('\n')
		if err != nil && !c.is_closed {
			fmt.Println("client: error ", err)
			return
		} else {
			if !c.is_closed {
				fmt.Println("client: recv ", res)
			}
		}
	}
}
