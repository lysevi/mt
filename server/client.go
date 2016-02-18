package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var _ = fmt.Sprintf("")

const (
	clientReadTimeOut = (time.Duration(300) * time.Millisecond)
)

type Client struct {
	name         string
	conn         net.Conn
	close_ch     chan interface{}
	is_closed    bool
	is_connected bool
	wg           *sync.WaitGroup
}

func Connect(name string, con_str string) (*Client, error) {
	conn, err := net.Dial("tcp", con_str)
	if err != nil {
		return nil, err
	}
	res := &Client{}
	res.name = name
	res.conn = conn
	res.is_closed = false
	res.is_connected = false
	res.close_ch = make(chan interface{})
	res.wg = &sync.WaitGroup{}

	res.wg.Add(1)
	go res.client_worker()
	return res, nil
}
func (c *Client) onClose() {
	c.is_connected = false
}
func (c *Client) Disconnect() {
	log.Println("client: disconnect...")
	c.conn.Write([]byte(disconnect))
	c.is_closed = true
	c.close_ch <- true
	c.wg.Wait()
}

func (c *Client) client_worker() {
	c.is_connected = true

	defer c.conn.Close()
	buf := make([]byte, 1024, 1024)
	protocol := NewClientProtocol(c)

L:
	for {
		select {
		case <-c.close_ch:
			{
				log.Println("client: stopChanel")
				c.onClose()
				break L
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
			log.Println("client: error ", c.is_closed, err)
			c.onClose()
			break L
		} else {
			if !c.is_closed && n != 0 {
				sb := string(buf[:n])
				log.Println("client: recv n: ", n, " buf:", strings.Replace(string(sb), "\n", "<", -1))
				protocol.OnRecv(buf[:n])

			}
		}
	}
	c.wg.Done()
	log.Println("client: done")
}

func (c *Client) Ping() {
	log.Println("client: ping")
	c.conn.Write([]byte(pong))
}

func (c *Client) SendName() {
	log.Println("client: send name")
	c.conn.Write([]byte(fmt.Sprintf("%s %s\n", helloFromClient, c.name)))
}

func (c *Client) Error(msg string) {
	log.Panicln(fmt.Sprint("server: error ", msg))
}
