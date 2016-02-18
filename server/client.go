package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

var _ = fmt.Sprintf("")

const (
	clientReadTimeOut = (time.Duration(300) * time.Millisecond)
)

type Client struct {
	conn         net.Conn
	close_ch     chan interface{}
	is_closed    bool
	is_connected bool
	wg           *sync.WaitGroup
}

func Connect(con_str string) (*Client, error) {
	conn, err := net.Dial("tcp", con_str)
	if err != nil {
		return nil, err
	}
	res := &Client{}
	res.conn = conn
	res.is_closed = false
	res.is_connected = false
	res.close_ch = make(chan interface{})
	res.wg = &sync.WaitGroup{}
	fmt.Println("add!")
	res.wg.Add(1)
	go res.client_worker()
	return res, nil
}
func (c *Client) onClose() {
	c.is_connected = false
}
func (c *Client) Disconnect() {
	c.is_closed = true
	c.close_ch <- true
	c.wg.Wait()
}

func (c *Client) client_worker() {
	c.is_connected = true

	defer c.conn.Close()
	buf := make([]byte, 1024, 1024)
	for {
		select {
		case <-c.close_ch:
			{
				fmt.Println("client: stopChanel")
				c.wg.Done()
				c.onClose()
				break
			}
		default:
		}
		c.conn.SetDeadline(time.Now().Add(clientReadTimeOut))

		n, err := c.conn.Read(buf)

		if err != nil && !c.is_closed {
			opErr, ok := err.(*net.OpError)
			if ok && (opErr.Timeout() || opErr.Err == io.EOF) {
				continue
			}
			fmt.Println("client: error ", c.is_closed, err)
			c.onClose()
			break
		} else {
			if !c.is_closed && n != 0 {
				fmt.Println("client: recv n: ", n, " buf:", string(buf))
			}
		}
	}
	fmt.Println("client done")
}
