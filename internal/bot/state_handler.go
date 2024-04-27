package bot

import (
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/events"
)

func (b *Bot) enterChannelName(userCache *cache.UserCache, msg events.Message) (bool, error) {
	b.Answer(msg).Text(msg.Ctx, msg.GetMessageText())
	userCache.State = cache.StateNone
	return true, nil
}

func (b *Bot) stateHandler(msg events.Message) (bool, error) {
	user := b.getOrCreateUser(msg.PeerUser.ID)
	if user.State == cache.StateNone {
		return false, nil
	}
	switch user.State {
	case cache.WaitingChannelName:
		return b.enterChannelName(user, msg)
	}

	return false, nil
}
