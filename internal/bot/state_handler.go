package bot

import (
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/events"
)

func (b *Bot) enterChannelName(msg events.MsgContext) (bool, error) {
	msg.UserCache.SetState(cache.StateNone)

	b.DeleteMessage(msg.Ctx, msg.Message.ID)

	channel, err := b.getChannelByName(msg.Ctx, msg.GetMessageText())
	if err != nil {
		b.Answer(msg.PeerUser).Text(msg.Ctx, "Ошибка при выполнении. Попробуйте позже.")
		return false, b.showMainPage(msg.Ctx, msg.PeerUser, msg.UserCache)
	}

	if msg.UserCache.HasChannelByID(channel.ID) {
		b.Answer(msg.PeerUser).NoWebpage().Textf(msg.Ctx, "Канал [%s](%s) уже был добавлен.", channel.Title, msg.GetMessageText())
		return false, b.showMainPage(msg.Ctx, msg.PeerUser, msg.UserCache)
	}

	_, err = msg.UserCache.AddGroup(channel)
	if err != nil {
		b.Answer(msg.PeerUser).Text(msg.Ctx, "Ошибка при выполнении. Попробуйте позже.")
		return false, b.showMainPage(msg.Ctx, msg.PeerUser, msg.UserCache)
	}

	b.Answer(msg.PeerUser).NoWebpage().Textf(msg.Ctx, "Канал [%s](%s) успешно добавлен.", channel.Title, msg.GetMessageText())
	return true, b.showChannelInfo(msg.Ctx, channel.ID, msg.PeerUser, msg.UserCache)
}

func (b *Bot) enterKeyWord(msg events.MsgContext) (bool, error) {
	msg.UserCache.SetState(cache.StateNone)

	b.DeleteMessage(msg.Ctx, msg.Message.ID)

	channel, ok := msg.UserCache.Channels[msg.UserCache.ActiveChannelID]
	if !ok {
		b.Answer(msg.PeerUser).Textf(msg.Ctx, "Ошибка при получении канала.")
		return false, b.showMainPage(msg.Ctx, msg.PeerUser, msg.UserCache)
	}

	if channel.AddKeyword(msg.GetMessageText()) != nil {
		b.Answer(msg.PeerUser).Textf(msg.Ctx, "Ключевое слово не было добавлено")
		return false, b.showChannelInfo(msg.Ctx, channel.TelegramID, msg.PeerUser, msg.UserCache)
	}

	return true, b.showChannelInfo(msg.Ctx, channel.TelegramID, msg.PeerUser, msg.UserCache)
}

func (b *Bot) stateHandler(msg events.MsgContext) (bool, error) {
	switch msg.UserCache.State {
	case cache.WaitingChannelName:
		return b.enterChannelName(msg)
	case cache.WaitingKeyWord:
		return b.enterKeyWord(msg)
	}

	return false, nil
}
