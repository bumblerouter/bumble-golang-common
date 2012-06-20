package message

import (
	"bumbleserver.org/common/peer"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Header struct {
	MessageType MessageType `json:"y"`
	MessageCode MessageCode `json:"c"`
	From        *peer.Peer  `json:"f,omitempty"`
	To          *peer.Peer  `json:"t,omitempty"`
	Date        time.Time   `json:"d,omitempty"`
}

func HeaderParse(s string) (*Header, error) {
	msg := new(Header)
	err := json.Unmarshal([]byte(s), msg)
	if err != nil {
		fmt.Printf("Header-PARSE-ERROR: %s\n", err)
		return nil, errors.New("unable to parse Header")
	}
	return msg, nil
}

func (p *Header) GetType() MessageType {
	return p.MessageType
}

func (p *Header) SetType(t MessageType) {
	p.MessageType = t
}

func (p *Header) GetCode() MessageCode {
	return p.MessageCode
}

func (p *Header) SetCode(c MessageCode) {
	p.MessageCode = c
}

func (p *Header) GetFrom() *peer.Peer {
	return p.From
}

func (p *Header) SetFrom(f *peer.Peer) {
	p.From = f
}

func (p *Header) GetTo() *peer.Peer {
	return p.To
}

func (p *Header) SetTo(t *peer.Peer) {
	p.To = t
}

func (p *Header) GetDate() time.Time {
	return p.Date
}

func (p *Header) SetDate(d time.Time) {
	p.Date = d
}
