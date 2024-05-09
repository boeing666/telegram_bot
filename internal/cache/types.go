package cache

import (
	"sync"

	"github.com/gotd/td/tg"
)

const (
	StateNone = iota
	WaitingChannelName
	WaitingKeyWord
)

type PeerKeyWords struct {
	Keywords map[int64]string
}

// Peer can be:
// tg.InputPeerChat
// tg.InputPeerChannel
type PeerInfo struct {
	TelegramID int64
	DatabaseID int64
	Name       string
	Title      string
	LastMsgID  int
	Peer       tg.InputPeerClass

	/* telegram user id -> keywords */
	UsersKeyWords map[int64]*PeerKeyWords
}

type UserData struct {
	Peer               tg.InputPeerUser
	State              uint32
	ActiveMenuID       int
	ActivePeerID       int64
	SecretButtonClicks int

	Channels map[int64]*PeerInfo
}

type PeersManager struct {
	Peers map[int64]*PeerInfo
	Users map[int64]*UserData
	Mutex sync.RWMutex
}
