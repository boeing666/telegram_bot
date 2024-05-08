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

type ChannelKeyWords struct {
	Keywords map[int64]string
}

type ChannelInfo struct {
	DatabaseID int64
	TelegramID int64
	Name       string
	Title      string
	LastMsgID  int
	Peer       tg.InputPeerClass

	/* telegram user id -> keywords */
	UsersKeyWords map[int64]*ChannelKeyWords
}

type UserData struct {
	TelegramID         int64
	State              uint32
	ActiveMenuID       int
	ActiveChannelID    int64
	SecretButtonClicks int
	Channels           map[int64]*ChannelInfo
}

type ChannelsManager struct {
	Channels map[int64]*ChannelInfo
	Users    map[int64]*UserData
	Mutex    sync.RWMutex
}

type rowGroups struct {
	ID        int64
	UserID    int64
	ChannelID int64
	LastMsgID int
	Name      string
	Title     string
}

type rowKeywords struct {
	DatabaseID int64
	GroupFK    int64
	Keyword    string
}

type UserKeyWords struct {
	userID  int64
	channel *ChannelInfo
}
