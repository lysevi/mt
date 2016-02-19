package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/lysevi/mt/storage"
)

var _ = fmt.Sprintf("")

const (
	clientReadTimeOut  = (time.Duration(300) * time.Millisecond)
	clientQueryTimeout = (time.Duration(3) * time.Second)
)

type Client struct {
	id           int32
	name         string
	conn         net.Conn
	conn_str     string
	close_ch     chan interface{}
	is_closed    bool
	is_connected bool
	wg           *sync.WaitGroup
	protocol     *ProtocolClient
}

func Connect(name string, con_str string) (*Client, error) {
	conn, err := net.Dial("tcp", con_str)
	if err != nil {
		return nil, err
	}
	res := &Client{}
	res.conn_str = con_str
	res.name = name
	res.conn = conn
	res.is_closed = false
	res.is_connected = false
	res.close_ch = make(chan interface{})
	res.wg = &sync.WaitGroup{}
	res.protocol = NewClientProtocol(res)

	buf, err := res.readnet()
	if err != nil {
		panic(err)
	}
	res.protocol.OnRecv(buf)
	if !res.is_connected {
		return nil, fmt.Errorf("connetion error. see log")
	}
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
	log.Println("client: disconnect ok.")
}

func (c *Client) readnet() ([]byte, error) {
	buf := make([]byte, 1024, 1024)
	c.conn.SetDeadline(time.Now().Add(clientReadTimeOut))
	n, err := c.conn.Read(buf)

	if err != nil && !c.is_closed {
		opErr, ok := err.(*net.OpError)
		if ok && (opErr.Timeout() || opErr.Err == io.EOF) {
			return nil, nil
		}
		log.Println("client: error ", c.is_closed, err)
		c.onClose()
		return nil, err
	} else {
		if !c.is_closed && n != 0 {
			sb := string(buf[:])
			log.Println("client: recv n: ", n, " buf:", strings.Replace(string(sb), "\n", "<", -1))
			return buf[:n], nil
		}
	}
	return nil, nil
}

func (c *Client) client_worker() {
	defer c.conn.Close()

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
		buf, e := c.readnet()
		if e != nil {
			panic(e)
		}
		if buf != nil {
			c.protocol.OnRecv(buf)
		}

	}
	c.wg.Done()
	log.Println("client: client_worker done")
}

func (c *Client) Ping() {
	log.Println("client: ping")
	c.conn.Write([]byte(pong))
}

func (c *Client) SendName() {
	//log.Println("client: send name")
	c.conn.Write([]byte(fmt.Sprintf("%s %s\n", helloFromClient, c.name)))
	buf := make([]byte, 100, 100)
	n, err := c.conn.Read(buf)
	if n == 0 || err != nil {
		panic(fmt.Sprintf("client: read id error %v n=%v", err, n))
	}

	n, err = fmt.Sscanf(string(buf), "%d", &c.id)
	if n == 0 || err != nil {
		panic(fmt.Sprintf("client: scan id error %v n=%v", err, n))
	}
	c.is_connected = true
	log.Println("client: connect ok")
}

func (c *Client) Error(msg string) {
	log.Panicln(fmt.Sprint("server: error ", msg))
}

func (c *Client) SendQuery(query []byte) ([]byte, error) {
	conn, err := net.Dial("tcp", c.conn_str)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 1024, 1024)
	n, err := conn.Read(buf)

	conn.SetReadDeadline(time.Now().Add(clientQueryTimeout))
	conn.Write([]byte(fmt.Sprintf("%s %d %s \n", queryRequest, c.id, string(query))))
	conn.Write([]byte(ok))

	result := []byte{}

L:
	for i := 0; ; i++ {
		if i > 1000 {
			break

		}
		log.Println("client: i:", i)
		buf := make([]byte, 1024, 1024)
		n, err = conn.Read(buf)
		log.Println("client: sendQuery recv ", string(buf[:n]))
		breader := bytes.NewBuffer(buf)
	SUB:
		for {
			bts, err := breader.ReadBytes('\n')
			if err == io.EOF {
				break SUB
			}
			log.Println("client: recv ", string(bts), len(bts), err)
			if IsError(bts) {
				panic(fmt.Sprintf("query error: ", string(bts[:n])))
			} else {
				if IsOk(bts) {
					log.Println("client: query end")
					break L
				} else {
					log.Println("client: query data  ", string(bts[:len(bts)-1]))
					result = append(result, bts...)
				}
			}
		}
	}
	log.Println("client: query result: ", len(result))
	return result, nil
}

func (c *Client) WriteValues(values []Value) error {
	qw := NewQueryWrite()
	qw.Values = values
	json_string, err := qw.bytes()
	if err != nil {
		panic(fmt.Sprintf("json error: %v", err))
	}
	//log.Println("client: ", string(json_string), err)

	_, err = c.SendQuery(json_string)
	return err
}

func (c *Client) ReadValues(from, to storage.Time) ([]Value, error) {
	qr := NewQueryRead()
	qr.From = from
	qr.To = to
	json_string, err := qr.bytes()
	if err != nil {
		panic(fmt.Sprintf("json error: %v", err))
	}
	//log.Println("client: ", string(json_string), err)

	answer, err := c.SendQuery(json_string)
	log.Println("client: answer ", string(answer))
	var values []Value
	err = json.Unmarshal(answer, &values)
	return values, err
}
