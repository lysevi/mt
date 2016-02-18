package server

import (
	"testing"
)

const (
	mok_action_ping      int = 1 << iota
	mok_action_pong      int = 1 << iota
	mok_action_say_hello int = 1 << iota
	mok_action_send_name int = 1 << iota
)

type mok_ServerActions struct {
	lastAction int
}
type mok_ClientActions struct {
	lastAction int
}

func (s *mok_ServerActions) Pong() {
	s.lastAction = mok_action_pong
}
func (s *mok_ServerActions) SayHello() {
	s.lastAction = mok_action_say_hello
}

func (s *mok_ClientActions) Ping() {
	s.lastAction = mok_action_ping
}

func (s *mok_ClientActions) SendName() {
	s.lastAction = mok_action_send_name
}

func TestProtocol(t *testing.T) {
	sa := mok_ServerActions{}
	ca := mok_ClientActions{}
	sp := NewServerProtocol(&sa)
	cp := NewClientProtocol(&ca)

	sp.OnRecv([]byte(helloFromClient))
	if sa.lastAction != mok_action_say_hello {
		t.Error(sa.lastAction, mok_action_say_hello)
	}

	cp.OnRecv([]byte(helloFromServer))
	if ca.lastAction != mok_action_send_name {
		t.Error(ca.lastAction, mok_action_send_name)
	}

	sp.OnRecv([]byte(pong))
	if sa.lastAction != mok_action_pong {
		t.Error(sa.lastAction, mok_action_pong)
	}

	cp.OnRecv([]byte(ping))
	if ca.lastAction != mok_action_ping {
		t.Error(ca.lastAction, mok_action_ping)
	}
}
