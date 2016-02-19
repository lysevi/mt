package server

import (
	"fmt"
	"strings"
)

const (
	helloFromClient = "***"
	helloFromServer = "***\n"
	disconnect      = "+Bye!\n"
	errorMsg        = "-\n"
	ping            = "+ping\n"
	pong            = "+pong\n"
	ok              = "+ok\n"
	queryRequest    = "+query"
)

func toBytes(s string) []byte {
	return []byte(s)
}

func IsOk(message []byte) bool {
	return strings.Compare(string(message), ok) == 0
}

func IsError(message []byte) bool {
	return len(message) > 0 && string(message)[0] == '-'
}

type ClientAction interface {
	Ping()
	SendName()
	Error(msg string) //server send info about error
}

type ServerAction interface {
	Pong(ci *ClientInfo) bool
	SayHello(ci *ClientInfo, buf []byte) bool
	Error(ci *ClientInfo, msg string) bool //client send info about error
	Disconnect(ci *ClientInfo) bool
	NewQuery(ci *ClientInfo, buf []byte) bool
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

func (p *ProtocolServer) OnRecv(ci *ClientInfo, message []byte) (bool, error) {
	hf_len := len(helloFromClient)
	if len(message) > hf_len && string(message[:hf_len]) == helloFromClient {
		return p.sa.SayHello(ci, message[hf_len+1:]), nil
	}

	if string(message) == pong {
		return p.sa.Pong(ci), nil
	}

	if len(message) > 0 && string(message)[0] == '-' {
		return p.sa.Error(ci, string(message)), nil
	}

	if string(message) == disconnect {
		return p.sa.Disconnect(ci), nil
	}

	if len(message) > len(queryRequest) && string(message[:len(queryRequest)]) == queryRequest {
		return p.sa.NewQuery(ci, message), nil
	}
	panic(fmt.Sprintf("ProtocolServer: uncknow command: '%v'", string(message)))
}

func NewClientProtocol(ca ClientAction) *ProtocolClient {
	res := ProtocolClient{}
	res.ca = ca
	return &res
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
