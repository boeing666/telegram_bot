package bot

import (
	"fmt"
	"tg_reader_bot/internal/cache"
	protos "tg_reader_bot/internal/protobufs"

	"github.com/gotd/td/tg"
	"google.golang.org/protobuf/proto"
)

func (b *Bot) callbackAddNewChannel(btn buttonContext) error {
	btn.UserCache.SetState(cache.WaitingChannelName)
	_, err := b.Answer(btn.User).Reply(btn.Update.MsgID).Text(btn.Ctx, "Введите в чат ссылку/айди имя чата/группы.")
	return err
}

func (b *Bot) callbackMyChannels(btn buttonContext) error {
	if len(btn.UserCache.Channels) == 0 {
		b.Answer(btn.User).Textf(btn.Ctx, "Вы не отслеживаете никакие каналы.")
		return nil
	}

	var rows []tg.KeyboardButtonRow
	for _, channel := range btn.UserCache.Channels {
		row := tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					fmt.Sprintf("%s (%d)", channel.Title, len(channel.KeyWords)),
					uint32(protos.MessageID_ChannelInfo),
					&protos.ButtonChanneInfo{Id: channel.TelegramID},
				),
			},
		}
		rows = append(rows, row)
	}
	rows = append(rows, CreateBackButton(uint32(protos.MessageID_MainPage)))

	_, err := b.Client.MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.Update.MsgID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     "Ваши отслеживаемые каналы, нажмите, чтобы настроить.",
	})

	return err
}

func (b *Bot) callbackChannelInfo(btn buttonContext) error {
	message := protos.ButtonChanneInfo{}
	proto.Unmarshal(btn.Data, &message)

	channel, ok := btn.UserCache.Channels[message.Id]
	if !ok {
		b.Answer(btn.User).Textf(btn.Ctx, "Ошибка при чтении ключевых слов.")
		return nil
	}

	rows := []tg.KeyboardButtonRow{
		{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					"Добавить новое",
					uint32(protos.MessageID_AddNewKeyWord),
					&protos.ButtonChanneInfo{Id: channel.TelegramID},
				),
			},
		},
	}

	for id, keyword := range channel.KeyWords {
		row := tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					keyword,
					uint32(protos.MessageID_RemoveKeyWord),
					&protos.ButtonRemoveKeyWord{KeywordId: id, GroupId: channel.TelegramID},
				),
			},
		}
		rows = append(rows, row)
	}
	rows = append(rows, CreateBackButton(uint32(protos.MessageID_MyChannels)))

	_, err := b.Client.MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.Update.MsgID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     fmt.Sprintf("Ключевые слова для канала %s, нажмите на слово, чтобы удалить.", channel.Name),
	})

	return err
}

func (b *Bot) callbackBack(btn buttonContext) error {
	message := protos.ButtonMenuBack{}
	proto.Unmarshal(btn.Data, &message)

	return b.btnCallbacks[message.Newmenu](btn)
}

func (b *Bot) callbackMainPage(btn buttonContext) error {
	_, err := b.Client.MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.Update.MsgID,
		ReplyMarkup: buildInitalMenu(),
		Message:     fmt.Sprintf("Добро пожаловать %s %s", btn.User.FirstName, btn.User.LastName),
	})
	return err
}

func (b *Bot) callbackAddNewKeyWord(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackRemoveKeyWord(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackNextChannels(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackPrevChannels(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackNextKeyWords(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackPrevKeyWords(btn buttonContext) error {
	return nil
}

func (b *Bot) registerQueryCallbacks() {
	b.btnCallbacks[uint32(protos.MessageID_AddNewChannel)] = b.callbackAddNewChannel
	b.btnCallbacks[uint32(protos.MessageID_MyChannels)] = b.callbackMyChannels
	b.btnCallbacks[uint32(protos.MessageID_AddNewKeyWord)] = b.callbackAddNewKeyWord
	b.btnCallbacks[uint32(protos.MessageID_RemoveKeyWord)] = b.callbackRemoveKeyWord
	b.btnCallbacks[uint32(protos.MessageID_NextChannels)] = b.callbackNextChannels
	b.btnCallbacks[uint32(protos.MessageID_PrevChannels)] = b.callbackNextKeyWords
	b.btnCallbacks[uint32(protos.MessageID_NextKeyWords)] = b.callbackPrevKeyWords
	b.btnCallbacks[uint32(protos.MessageID_Back)] = b.callbackBack
	b.btnCallbacks[uint32(protos.MessageID_MainPage)] = b.callbackMainPage
	b.btnCallbacks[uint32(protos.MessageID_ChannelInfo)] = b.callbackChannelInfo
}
