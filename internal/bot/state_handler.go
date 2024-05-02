package bot

import (
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/events"

	"github.com/gotd/td/telegram/message/peer"
	"github.com/gotd/td/tg"
)

func (b *Bot) enterChannelName(msg events.MsgContext) (bool, error) {
	msg.UserCache.State = cache.StateNone

	otherPeer, err := b.Sender.Resolve(msg.GetMessageText()).AsInputPeer(msg.Ctx)
	if err != nil {
		b.Answer(msg).Text(msg.Ctx, "Введены некорректные данные.")
		return false, nil
	}

	switch otherPeer.(type) {
	case *tg.InputPeerChat:
		// chatPeer := peer.(*tg.InputPeerChat)
	case *tg.InputPeerChannel:
		inputChannel, ok := peer.ToInputChannel(otherPeer)
		if !ok {
			b.Answer(msg).Text(msg.Ctx, "Внутрення ошибка 101. Попробуйте позже.")
			return false, nil
		}
		channels, err := b.Client.ChannelsGetChannels(msg.Ctx, []tg.InputChannelClass{inputChannel})
		if err != nil {
			b.Answer(msg).Text(msg.Ctx, "Внутрення ошибка 102. Попробуйте позже.")
			return false, nil
		}
		channel := channels.GetChats()[0].(*tg.Channel)
		if msg.UserCache.HasChannelByID(channel.ID) {
			b.Answer(msg).Textf(msg.Ctx, "Канал [%s](%s) уже добавлен.", channel.Title, msg.GetMessageText())
			return false, nil
		}

		b.Answer(msg).Textf(msg.Ctx, "Канал [%s](%s) успешно добавлен.", channel.Title, msg.GetMessageText())
	default:
		b.Answer(msg).Textf(msg.Ctx, "[%s] не является каналом или чатом.", msg.GetMessageText())
		return false, nil
	}

	return true, nil
}

func (b *Bot) stateHandler(msg events.MsgContext) (bool, error) {
	switch msg.UserCache.State {
	case cache.WaitingChannelName:
		return b.enterChannelName(msg)
	}

	return false, nil
}
