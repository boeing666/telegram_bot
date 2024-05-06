package bot

import (
	"fmt"
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/events"
	"tg_reader_bot/internal/protobufs"

	"google.golang.org/protobuf/proto"
)

func (b *Bot) enterChannelName(msg events.MsgContext) (bool, error) {
	msg.UserCache.SetState(cache.StateNone)

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

func (b *Bot) enterKeyWord(msg events.MsgContext) (bool, error) {
	msg.UserCache.SetState(cache.StateNone)

	channel, ok := msg.UserCache.Channels[msg.UserCache.ActiveChannelID]
	if !ok {
		b.Answer(msg.PeerUser).Textf(msg.Ctx, "Ошибка при получении канала.")
		return false, nil
	}

	if channel.AddKeyword(msg.GetMessageText()) != nil {
		b.Answer(msg.PeerUser).Textf(msg.Ctx, "Ключевое слово не было добавлено")
		return false, nil
	}

	message := protobufs.ButtonChanneInfo{Id: channel.TelegramID}
	data, _ := proto.Marshal(&message)
	err := b.showChannelInfo(msg.Ctx, data, msg.PeerUser, msg.UserCache)
	if err != nil {
		return false, err
	}

	return true, nil
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
