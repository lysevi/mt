package server

import (
	"testing"
)

const (
	mok_action_ping      int = 1 << iota
	mok_action_pong      int = 1 << iota
	mok_action_say_hello int = 1 << iota
	mok_action_send_name int = 1 << iota
	mok_action_error     int = 1 << iota
	mok_action_discon    int = 1 << iota
	mok_action_newquery  int = 1 << iota
)

type mok_ServerActions struct {
	lastAction int
	msg        string
}
type mok_ClientActions struct {
	lastAction int
	msg        string
}

func (s *mok_ServerActions) Pong(ci *ClientInfo) bool {
	s.lastAction = mok_action_pong
	return false
}
func (s *mok_ServerActions) SayHello(ci *ClientInfo, buf []byte) bool {
	s.lastAction = mok_action_say_hello
	s.msg = string(buf)
	return false
}

func (s *mok_ServerActions) Error(ci *ClientInfo, msg string) bool {
	s.lastAction = mok_action_error
	s.msg = msg
	return false
}

func (s *mok_ServerActions) Disconnect(ci *ClientInfo) bool {
	s.lastAction = mok_action_discon
	return false
}

func (s *mok_ServerActions) NewQuery(ci *ClientInfo, buf []byte) bool {
	s.lastAction = mok_action_newquery
	return false
}

func (s *mok_ClientActions) Ping() {
	s.lastAction = mok_action_ping
}

func (s *mok_ClientActions) SendName() {
	s.lastAction = mok_action_send_name
}
func (s *mok_ClientActions) Error(msg string) {
	s.lastAction = mok_action_error
	s.msg = msg
}

func TestProtocol(t *testing.T) {
	sa := mok_ServerActions{}
	ca := mok_ClientActions{}
	sp := NewServerProtocol(&sa)
	cp := NewClientProtocol(&ca)

	sp.OnRecv(nil, []byte(helloFromClient+" test"))
	if sa.lastAction != mok_action_say_hello {
		t.Error(sa.lastAction, mok_action_say_hello)
	}

	cp.OnRecv([]byte(helloFromServer))
	if ca.lastAction != mok_action_send_name {
		t.Error(ca.lastAction, mok_action_send_name)
	}

	sp.OnRecv(nil, []byte(pong))
	if sa.lastAction != mok_action_pong {
		t.Error(sa.lastAction, mok_action_pong)
	}

	cp.OnRecv([]byte(ping))
	if ca.lastAction != mok_action_ping {
		t.Error(ca.lastAction, mok_action_ping)
	}

	cp.OnRecv([]byte(errorMsg))
	if ca.lastAction != mok_action_error || ca.msg != errorMsg {
		t.Error(ca.lastAction, mok_action_error, sa.msg)
	}

	sp.OnRecv(nil, []byte(errorMsg))
	if sa.lastAction != mok_action_error || sa.msg != errorMsg {
		t.Error(sa.lastAction, mok_action_error, sa.msg)
	}

	sp.OnRecv(nil, []byte(queryRequest+" qwerty"))
	if sa.lastAction != mok_action_newquery {
		t.Error(sa.lastAction, mok_action_error)
	}
}
