package bot

import (
	"time"
)

const (
	StateNone = iota
	WaitingChannelName
	WaitingKeyWord
)

type userData struct {
	state uint32
}

func createUser() userData {
	return userData{state: StateNone}
}

func (b *Bot) getOrCreateUser(userID int64) *userData {
	res, err := b.cache.Value(userID)
	if err == nil {
		return res.Data().(*userData)
	}
	user := createUser()
	b.cache.Add(userID, 10*time.Minute, &user)
	return &user
}

func (b *Bot) setUserState(userID int64, state uint32) {
	user := b.getOrCreateUser(userID)
	user.state = state
}

func (b *Bot) getUserState(userID int64) uint32 {
	user := b.getOrCreateUser(userID)
	return user.state
}
