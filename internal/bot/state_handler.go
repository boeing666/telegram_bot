package bot

import (
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/events"
)

func (b *Bot) enterChannelName(msg events.MsgContext) error {
	user := msg.UserData
	user.State = cache.StateNone

	b.Sender.Delete().Messages(msg.Ctx, msg.Message.ID)

	channel, err := GetChannelByName(b.API(), b.Sender, msg.Ctx, msg.GetText())
	if err != nil {
		b.Answer(msg.PeerUser).Text(msg.Ctx, "Ошибка при выполнении. Попробуйте позже.")
		return err
	}

	if user.HasChannelByID(channel.ID) {
		b.Answer(msg.PeerUser).NoWebpage().Textf(msg.Ctx, "Канал [%s](%s) уже был добавлен.", channel.Title, msg.GetText())
		return err
	}

	channelName := channel.Username
	if len(channelName) == 0 {
		channelName = channel.Usernames[0].Username
	}

	channelInfo, err := b.channelsCache.AddChannelToUser(user.GetID(), 0, channel.ID, 0, channelName, channel.Title, true)
	if err != nil {
		b.Answer(msg.PeerUser).Text(msg.Ctx, "Ошибка при выполнении. Попробуйте позже.")
		return err
	}
	channelInfo.Peer = channel.AsInputPeer()

	b.Answer(msg.PeerUser).NoWebpage().Textf(msg.Ctx, "Канал [%s](%s) успешно добавлен.", channel.Title, msg.GetText())
	return b.showChannelInfo(msg.Ctx, channel.ID, msg.PeerUser, user)
}

func (b *Bot) enterKeyWord(msg events.MsgContext) error {
	user := msg.UserData
	user.State = cache.StateNone

	b.Sender.Delete().Messages(msg.Ctx, msg.Message.ID)

	channel := user.GetActiveChannel()
	if channel == nil {
		b.Answer(msg.PeerUser).Textf(msg.Ctx, "Ошибка при получении канала.")
		return b.showMainPage(msg.Ctx, msg.PeerUser, user)
	}

	err := channel.AddKeyword(user.GetID(), 0, msg.GetText(), true)
	if err != nil {
		return err
	}

	return b.showChannelInfo(msg.Ctx, user.GetActivePeerID(), msg.PeerUser, user)
}

func (b *Bot) stateHandler(msg events.MsgContext) (bool, error) {
	switch msg.UserData.State {
	case cache.WaitingChannelName:
		return true, b.enterChannelName(msg)
	case cache.WaitingKeyWord:
		return true, b.enterKeyWord(msg)
	}

	return false, nil
}
