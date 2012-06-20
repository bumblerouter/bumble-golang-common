package message

import (
	"bumbleserver.org/common/peer"
	"time"
)

type MessageType int
type MessageCode int

const (
	CODE_GENERICERROR         MessageCode = 0
	CODE_AUTHENTICATE         MessageCode = 1
	CODE_AUTHENTICATION       MessageCode = 2
	CODE_AUTHENTICATIONRESULT MessageCode = 3
	CODE_GENERICMESSAGE       MessageCode = 4
)

type Message interface {
	GetType() MessageType
	SetType(MessageType)
	GetCode() MessageCode
	SetCode(MessageCode)
	GetFrom() *peer.Peer
	SetFrom(*peer.Peer)
	GetTo() *peer.Peer
	SetTo(*peer.Peer)
	GetDate() time.Time
	SetDate(time.Time)
}
