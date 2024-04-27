package bot

import (
	"tg_reader_bot/internal/cache"
	"time"
)

func (b *Bot) getOrCreateUser(userID int64) *cache.UserCache {
	res, err := b.cache.Value(userID)
	if err == nil {
		return res.Data().(*cache.UserCache)
	}
	user := cache.CreateUser()
	b.cache.Add(userID, 10*time.Minute, &user)
	return &user
}

func (b *Bot) setUserState(userID int64, state uint32) {
	user := b.getOrCreateUser(userID)
	user.State = state
}

func (b *Bot) getUserState(userID int64) uint32 {
	user := b.getOrCreateUser(userID)
	return user.State
}
