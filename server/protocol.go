package server

import (
	"reflect"
)

const (
	helloFromClient = "+++"
	helloFromServer = "+++"

	ping = "ping"
	pong = "pong"
)

func toBytes(s string) []byte {
	return []byte(s)
}

type ClientAction interface {
	Ping()
	SendName()
}

type ServerAction interface {
	Pong()
	SayHello()
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
	return nil
}
