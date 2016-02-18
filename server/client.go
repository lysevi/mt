package server

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var _ = fmt.Sprintf("")

type Client struct {
	conn         net.Conn
	is_closed    bool
	is_connected bool
	wg           sync.WaitGroup
}

func Connect(con_str string) (*Client, error) {
	conn, err := net.Dial("tcp", con_str)
	if err != nil {
		return nil, err
	}
	res := Client{}
	res.conn = conn
	res.is_closed = false
	res.is_connected = false
	res.wg.Add(1)
	go res.client_worker()
	return &res, nil
}

func (c *Client) Disconnect() {
	c.is_closed = true
	c.conn.Close()
	c.wg.Wait()
}

func (c *Client) client_worker() {
	reader := bufio.NewReader(c.conn)
	for {
		res, err := reader.ReadString('\n')
		c.is_connected = true
		if err != nil && !c.is_closed {
			fmt.Println("client: error ", err)
			break
		} else {
			if !c.is_closed {
				fmt.Println("client: recv ", res)
			} else {
				break
			}
		}
	}
	c.is_connected = false
	c.wg.Done()
}
