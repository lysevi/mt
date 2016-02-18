package server

import (
	"fmt"
	"strings"
)

const (
	helloFromClient = "+**"
	helloFromServer = "+++\n"
	disconnect      = "+Bye!\n"
	errorMsg        = "-\n"
	ping            = "+ping\n"
	pong            = "+pong\n"
	ok              = "+ok"
	query           = "+query"
)

func toBytes(s string) []byte {
	return []byte(s)
}

type ClientAction interface {
	Ping()
	SendName()
	Error(msg string) //server send info about error
}

type ServerAction interface {
	Pong(ci *ClientInfo)
	SayHello(ci *ClientInfo, buf []byte)
	Error(ci *ClientInfo, msg string) //client send info about error
	Disconnect(ci *ClientInfo)
}

type ProtocolServer struct {
	sa ServerAction
}

type ProtocolClient struct {
	ca ClientAction
}

func NewServerProtocol(sa ServerAction) ProtocolServer {
	res := ProtocolServer{}
	res.sa = sa
	return res
}

func (p *ProtocolServer) OnRecv(ci *ClientInfo, message []byte) error {
	hf_len := len(helloFromClient)
	if len(message) > hf_len && string(message[:hf_len]) == helloFromClient {
		p.sa.SayHello(ci, message[hf_len+1:])
		return nil
	}

	if string(message) == pong {
		p.sa.Pong(ci)
		return nil
	}

	if len(message) > 0 && string(message)[0] == '-' {
		p.sa.Error(ci, string(message))
		return nil
	}

	if string(message) == disconnect {
		p.sa.Disconnect(ci)
		return nil
	}
	panic(fmt.Sprintf("ProtocolServer: uncknow command: '%v'", string(message)))
}

func NewClientProtocol(ca ClientAction) ProtocolClient {
	res := ProtocolClient{}
	res.ca = ca
	return res
}

func (p *ProtocolClient) OnRecv(message []byte) error {
	if string(message) == helloFromServer {
		p.ca.SendName()
		return nil
	}

	if strings.Compare(string(message), ping) == 0 {
		p.ca.Ping()
		return nil
	}

	if len(message) > 0 && string(message)[0] == '-' {
		p.ca.Error(string(message))
		return nil
	}

	panic(fmt.Sprint("ProtocolClient: unknow command: ", string(message)))
}
