package bot

import (
	"tg_reader_bot/internal/cache"
	"time"
)

func (b *Bot) getOrCreateUser(userID int64, create bool) *cache.UserCache {
	res, err := b.cache.Value(userID)
	if err == nil {
		return res.Data().(*cache.UserCache)
	}
	if !create {
		return nil
	}
	user := cache.CreateUser(userID)
	b.cache.Add(userID, 10*time.Minute, user)
	return user
}
