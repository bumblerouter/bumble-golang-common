package session

import (
	"bumbleserver.org/common/envelope"
	"bumbleserver.org/common/peer"
	"code.google.com/p/go.net/websocket"
	"errors"
	"fmt"
	"io"
	"syscall"
	"time"
)

var sessionTimeout int = 60 // 60 seconds after disconnection the session will be killed

type SessionStore map[string]*Session

var sessionStore SessionStore = make(map[string]*Session)

type Session struct {
	peer       *peer.Peer
	websocket  *websocket.Conn
	queue      *envelope.Queue
	timeout    *time.Timer
	reset      chan bool
	connected  bool
	disconnect chan bool
}

func Config(st int) {
	sessionTimeout = st
}

func NewSession(p *peer.Peer, ws *websocket.Conn) *Session {
	if p.String() == "" {
		fmt.Printf("NEWSESSION CALLED ON %s but it was empty.\n", p)
		return nil
	}
	var s *Session
	if _, ok := sessionStore[p.String()]; ok {
		if sessionStore[p.String()].connected {
			fmt.Printf("NEWSESSION CALLED ON %s but it is already in use.\n", p)
			return nil
		} else {
			fmt.Printf("NEWSESSION CALLED ON %s so we reattached the existing one.\n", p)
			s = sessionStore[p.String()]
			s.peer = p
			s.websocket = ws
			s.timeout.Stop()
			s.connected = true
			s.reset <- true
		}
	} else {
		fmt.Printf("NEWSESSION CALLED ON %s so we made a new one.\n", p)
		s = new(Session)
		s.peer = p
		s.queue = envelope.NewQueue()
		s.websocket = ws
		s.timeout = new(time.Timer)
		s.reset = make(chan bool)
		s.disconnect = make(chan bool)
		s.connected = true
		sessionStore[p.String()] = s
		go s.queueRunner()
	}
	return s
}

func EndSession(p *peer.Peer) {
	if p.String() == "" {
		return
	}
	if _, ok := sessionStore[p.String()]; ok {
		// if there are messages in the queue, do we destroy them or return to sender or ...?  FIXME
		delete(sessionStore, p.String())
		fmt.Printf("ENDSESSION CALLED ON %s and so it was ended.\n", p)
		return
	}
	fmt.Printf("ENDSESSION CALLED ON %s but no session exists for that name.\n", p)
}

func DisconnectSession(p *peer.Peer) {
	if p.String() == "" {
		return
	}
	if _, ok := sessionStore[p.String()]; ok {
		fmt.Printf("DISCONNECTSESSION CALLED ON %s and so it was signaled.\n", p)
		sessionStore[p.String()].disconnect <- true
		return
	}
	fmt.Printf("DISCONNECTSESSION CALLED ON %s but no session exists for that name.\n", p)
}

func IsConnectedSession(p *peer.Peer) bool {
	if p.String() != "" {
		if _, ok := sessionStore[p.String()]; ok {
			return sessionStore[p.String()].connected
		}
	}
	return false
}

func PassEnvelope(p *peer.Peer, e *envelope.Envelope) error {
	if p.String() != "" {
		if _, ok := sessionStore[p.String()]; ok {
			sessionStore[p.String()].queue.Add(e)
			return nil
		}
	}
	return errors.New("requested peer session does not exist")
}

func (s *Session) EndSession() {
	EndSession(s.peer)
}

func (s *Session) DisconnectSession() {
	DisconnectSession(s.peer)
}

func (s *Session) IsConnectedSession() bool {
	return IsConnectedSession(s.peer)
}

func (s *Session) queueRunner() {
	// fmt.Println("SQR ENTRY")
	// defer fmt.Println("SQR EXIT")
	defer s.EndSession()
	for {
		if s.IsConnectedSession() {
			queue := s.queue.Channel()
			select {
			case <-s.reset: // we probably reconnected
				continue
			case <-s.disconnect: // we got disconnected
				s.connected = false
				continue
			case e := <-queue:
				err := websocket.JSON.Send(s.websocket, e)
				if err != nil {
					fmt.Printf("SQR-JSON-SEND ERROR: %s\n", err.Error())
					continue
				}
				if err == io.EOF || err == syscall.EINVAL { // peer disconnected (FIXME: want to get proper test for "read tcp ... use of closed network connection" error)
					queue <- e // put the genie back in the bottle
					s.connected = false
					continue
				}
			}
		} else { // we're not connected, so we want to wait for a timeout
			s.timeout = time.NewTimer(time.Duration(int64(sessionTimeout) * 1e9)) // set to a reasonable timeout period
			select {
			case <-s.reset: // we probably reconnected
				s.connected = true
				continue
			case <-s.disconnect: // we got disconnected
				continue
			case <-s.timeout.C:
				// fmt.Println("SQR DISCONNECTION TIMEOUT REACHED")
				return
			}
		}
	}
}
