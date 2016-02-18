package server

import (
	"reflect"
)

const (
	helloFromClient = "+++"
	helloFromServer = "+++"
	errorMsg        = "-"
	ping            = "ping"
	pong            = "pong"
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
	Pong()
	SayHello()
	Error(msg string) //client send info about error
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

func (p *ProtocolServer) OnRecv(message []byte) error {
	if reflect.DeepEqual(message, toBytes(helloFromClient)) {
		p.sa.SayHello()
	}

	if reflect.DeepEqual(message, toBytes(pong)) {
		p.sa.Pong()
	}

	if len(message) > 0 && string(message)[0] == '-' {
		p.sa.Error(string(message))
	}
	return nil
}

func NewClientProtocol(ca ClientAction) ProtocolClient {
	res := ProtocolClient{}
	res.ca = ca
	return res
}

func (p *ProtocolClient) OnRecv(message []byte) error {
	if reflect.DeepEqual(message, toBytes(helloFromServer)) {
		p.ca.SendName()
	}

	if reflect.DeepEqual(message, toBytes(ping)) {
		p.ca.Ping()
	}

	if len(message) > 0 && string(message)[0] == '-' {
		p.ca.Error(string(message))
	}
	return nil
}
