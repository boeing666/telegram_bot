package bot

import (
	"fmt"
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/events"
)

func (b *Bot) enterChannelName(msg events.MsgContext) (bool, error) {
	msg.UserCache.State = cache.StateNone

	channel, err := b.getChannelByName(msg.Ctx, msg.GetMessageText())
	if err != nil {
		fmt.Println("GetChannelByName", err)
		b.Answer(msg.PeerUser).Text(msg.Ctx, "Ошибка при выполнении. Попробуйте позже.")
		return false, nil
	}

	if msg.UserCache.HasChannelByID(channel.ID) {
		b.Answer(msg.PeerUser).Textf(msg.Ctx, "Канал [%s](%s) уже был добавлен.", channel.Title, msg.GetMessageText())
		return false, nil
	}

	_, err = msg.UserCache.AddGroup(channel)
	if err != nil {
		fmt.Println("AddGroup", err)
		b.Answer(msg.PeerUser).Text(msg.Ctx, "Ошибка при выполнении. Попробуйте позже.")
		return false, nil
	}

	b.Answer(msg.PeerUser).Textf(msg.Ctx, "Канал [%s](%s) успешно добавлен.", channel.Title, msg.GetMessageText())
	return true, nil
}

func (b *Bot) stateHandler(msg events.MsgContext) (bool, error) {
	switch msg.UserCache.State {
	case cache.WaitingChannelName:
		return b.enterChannelName(msg)
	}

	return false, nil
}
