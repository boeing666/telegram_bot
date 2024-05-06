package bot

import (
	"context"
	"fmt"
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/protobufs"

	"github.com/gotd/td/tg"
	"google.golang.org/protobuf/proto"
)

func (b *Bot) callbackAddNewChannel(btn buttonContext) error {
	rows := []tg.KeyboardButtonRow{CreateBackButton("Отмена", protobufs.MessageID_MainPage, nil)}
	btn.UserCache.SetState(cache.WaitingChannelName)
	_, err := b.Client.MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserCache.ActiveMenuID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     "Введите в чат ссылку/айди имя чата/группы.",
	})
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
					protobufs.MessageID_ChannelInfo,
					&protobufs.ButtonChanneInfo{Id: channel.TelegramID},
				),
			},
		}
		rows = append(rows, row)
	}
	rows = append(rows, CreateBackButton("Назад", protobufs.MessageID_MainPage, nil))

	_, err := b.Client.MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserCache.ActiveMenuID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     "Ваши отслеживаемые каналы, нажмите, чтобы настроить.",
	})

	return err
}

func (b *Bot) showChannelInfo(ctx context.Context, data []byte, User *tg.User, userCache *cache.UserCache) error {
	var message protobufs.ButtonChanneInfo
	proto.Unmarshal(data, &message)

	channel, ok := userCache.Channels[message.Id]
	if !ok {
		b.Answer(User).Textf(ctx, "Ошибка при чтении ключевых слов.")
		return nil
	}

	userCache.ActiveChannelID = channel.TelegramID

	rows := []tg.KeyboardButtonRow{
		{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					"Добавить новое",
					protobufs.MessageID_AddNewKeyWord,
					&protobufs.ButtonChanneInfo{Id: channel.TelegramID},
				),
			},
		},
	}

	for id, keyword := range channel.KeyWords {
		row := tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					keyword,
					protobufs.MessageID_RemoveKeyWord,
					&protobufs.ButtonRemoveKeyWord{KeywordId: id, GroupId: channel.TelegramID},
				),
			},
		}
		rows = append(rows, row)
	}
	rows = append(rows,
		CreateBackButton("Назад", protobufs.MessageID_MyChannels, nil),
		CreateBackButton("На главную", protobufs.MessageID_MainPage, nil),
	)

	_, err := b.Client.MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: userCache.TelegramID},
		ID:          userCache.ActiveMenuID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     fmt.Sprintf("Ключевые слова для канала %s, нажмите на слово, чтобы удалить.", channel.Name),
	})

	return err
}

func (b *Bot) callbackChannelInfo(btn buttonContext) error {
	return b.showChannelInfo(btn.Ctx, btn.Data, btn.User, btn.UserCache)
}

func (b *Bot) callbackBack(btn buttonContext) error {
	message := protobufs.ButtonMenuBack{}
	proto.Unmarshal(btn.Data, &message)

	if message.Msg != nil {
		btn.Data = message.Msg
	}

	return b.btnCallbacks[message.Newmenu](btn)
}

func (b *Bot) callbackMainPage(btn buttonContext) error {
	_, err := b.Client.MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserCache.ActiveMenuID,
		ReplyMarkup: buildInitalMenu(),
		Message:     fmt.Sprintf("Добро пожаловать %s %s", btn.User.FirstName, btn.User.LastName),
	})
	return err
}

func (b *Bot) callbackAddNewKeyWord(btn buttonContext) error {
	var message protobufs.ButtonChanneInfo
	proto.Unmarshal(btn.Data, &message)

	channel, ok := btn.UserCache.Channels[message.Id]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при чтении ключевых слов.")
		return nil
	}

	btn.UserCache.SetState(cache.WaitingKeyWord)

	rows := []tg.KeyboardButtonRow{CreateBackButton("Отмена", protobufs.MessageID_ChannelInfo, &message)}
	_, err := b.Client.MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserCache.ActiveMenuID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     fmt.Sprintf("Канал: [%s]\nВведите ключевое слово или регулярное выражение.", channel.Title),
	})

	return err
}

func (b *Bot) callbackRemoveKeyWord(btn buttonContext) error {
	var message protobufs.ButtonRemoveKeyWord
	proto.Unmarshal(btn.Data, &message)

	channel, ok := btn.UserCache.Channels[message.GroupId]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при удалении.")
		return nil
	}

	err := channel.RemoveKeyword(message.KeywordId)
	if err != nil {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при удалении.")
		return nil
	}

	return b.callbackChannelInfo(btn)
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
	b.btnCallbacks[protobufs.MessageID_AddNewChannel] = b.callbackAddNewChannel
	b.btnCallbacks[protobufs.MessageID_MyChannels] = b.callbackMyChannels
	b.btnCallbacks[protobufs.MessageID_AddNewKeyWord] = b.callbackAddNewKeyWord
	b.btnCallbacks[protobufs.MessageID_RemoveKeyWord] = b.callbackRemoveKeyWord
	b.btnCallbacks[protobufs.MessageID_NextChannels] = b.callbackNextChannels
	b.btnCallbacks[protobufs.MessageID_PrevChannels] = b.callbackNextKeyWords
	b.btnCallbacks[protobufs.MessageID_NextKeyWords] = b.callbackPrevKeyWords
	b.btnCallbacks[protobufs.MessageID_Back] = b.callbackBack
	b.btnCallbacks[protobufs.MessageID_MainPage] = b.callbackMainPage
	b.btnCallbacks[protobufs.MessageID_ChannelInfo] = b.callbackChannelInfo
}
