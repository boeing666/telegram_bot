package bot

import (
	"context"
	"fmt"
	"tg_reader_bot/internal/cache"
	"time"

	"github.com/gotd/td/tg"
)

func (b *Bot) getOrCreateUser(ctx context.Context, peerUser *tg.User, create bool) (*cache.UserCache, error) {
	res, err := b.cache.Value(peerUser.ID)
	if err == nil {
		return res.Data().(*cache.UserCache), nil
	}

	if !create {
		return nil, nil
	}

	user, err := cache.CreateUser(peerUser)
	b.updateUserGroupsInfo(ctx, user)
	b.cache.Add(peerUser.ID, 10*time.Minute, user)

	return user, err
}

func (b *Bot) updateUserGroupsInfo(ctx context.Context, userCache *cache.UserCache) {
	for _, channel := range userCache.Channels {
		peer, err := b.getChannelByName(ctx, channel.Name)
		if err == nil {
			channel.Title = peer.Title
		} else {
			fmt.Printf("Error resolve %v | %v\n", channel, err)
		}
	}
}
