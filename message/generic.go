package message

import (
	"bumbleserver.org/common/peer"
	"encoding/json"
	"errors"
	"time"
)

const TYPE_GENERIC MessageType = 0

type Generic struct {
	MessageType MessageType `json:"y"`
	MessageCode MessageCode `json:"c"`
	From        *peer.Peer  `json:"f,omitempty"`
	To          *peer.Peer  `json:"t,omitempty"`
	Date        time.Time   `json:"d,omitempty"`
	Error       string      `json:"e,omitempty"`
	Info        string      `json:"i,omitempty"`
	Success     bool        `json:"b,omitempty"`
}

func NewGeneric(c MessageCode) *Generic {
	p := new(Generic)
	p.SetType(TYPE_GENERIC)
	p.SetCode(c)
	return p
}

func GenericParse(s string) (*Generic, error) {
	msg := new(Generic)
	err := json.Unmarshal([]byte(s), msg)
	if err != nil {
		return nil, errors.New("unable to parse message.Generic")
	}
	return msg, nil
}

func (p *Generic) GetType() MessageType {
	return p.MessageType
}

func (p *Generic) SetType(t MessageType) {
	p.MessageType = t
}

func (p *Generic) GetCode() MessageCode {
	return p.MessageCode
}

func (p *Generic) SetCode(c MessageCode) {
	p.MessageCode = c
}

func (p *Generic) GetFrom() *peer.Peer {
	return p.From
}

func (p *Generic) SetFrom(f *peer.Peer) {
	p.From = f
}

func (p *Generic) GetTo() *peer.Peer {
	return p.To
}

func (p *Generic) SetTo(t *peer.Peer) {
	p.To = t
}

func (p *Generic) GetDate() time.Time {
	return p.Date
}

func (p *Generic) SetDate(d time.Time) {
	p.Date = d
}

func (p *Generic) GetError() string {
	return p.Error
}

func (p *Generic) SetError(e string) {
	p.Error = e
}

func (p *Generic) GetInfo() string {
	return p.Info
}

func (p *Generic) SetInfo(i string) {
	p.Info = i
}

func (p *Generic) GetSuccess() bool {
	return p.Success
}

func (p *Generic) SetSuccess(s bool) {
	p.Success = s
}
