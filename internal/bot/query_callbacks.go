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
	btn.UserData.State = cache.WaitingChannelName
	_, err := b.API().MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserData.ActiveMenuID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     "Введите в чат ссылку/айди имя чата/группы.",
	})
	return err
}

func (b *Bot) callbackMyChannels(btn buttonContext) error {
	if len(btn.UserData.Channels) == 0 {
		return b.SetAnswerCallback(btn.Ctx, "Список каналов пуст", btn.Update.QueryID)
	}

	var rows []tg.KeyboardButtonRow
	for _, channel := range btn.UserData.Channels {
		row := tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					fmt.Sprintf("%s (%d)", channel.Title, channel.GetUserKeyWordsCount(btn.UserData.TelegramID)),
					protobufs.MessageID_ChannelInfo,
					&protobufs.ButtonChanneInfo{ChannelId: channel.TelegramID},
				),
			},
		}
		rows = append(rows, row)
	}
	rows = append(rows,
		CreateSpaceButtonRow(),
		CreateBackButton("Назад", protobufs.MessageID_MainPage, nil),
	)

	_, err := b.API().MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserData.ActiveMenuID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     "Ваши отслеживаемые каналы, нажмите, чтобы настроить.\n ",
	})

	return err
}

func (b *Bot) showChannelInfo(ctx context.Context, channelId int64, User *tg.User, user *cache.UserData) error {
	channel, ok := user.Channels[channelId]
	if !ok {
		b.Answer(User).Textf(ctx, "Ошибка при поиске канала.")
		return nil
	}

	user.ActiveChannelID = channel.TelegramID

	rows := []tg.KeyboardButtonRow{
		{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					"Добавить новое ключевое слово",
					protobufs.MessageID_AddNewKeyWord,
					&protobufs.ButtonChanneInfo{ChannelId: channel.TelegramID},
				),
			},
		},
	}

	keywords := channel.GetUserKeyWords(user.TelegramID)
	if keywords != nil {
		for id, keyword := range *keywords {
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
	}

	rows = append(rows,
		CreateSpaceButtonRow(),
		CreateButtonRow("Удалить канал", protobufs.MessageID_RemoveChannel, &protobufs.ButtonChanneInfo{ChannelId: channel.TelegramID}),
		CreateBackButton("Назад", protobufs.MessageID_MyChannels, nil),
		CreateBackButton("На главную", protobufs.MessageID_MainPage, nil),
	)

	_, err := b.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: user.TelegramID},
		ID:          user.ActiveMenuID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     fmt.Sprintf("Ключевые слова для канала %s(%s)\nНажмите на слово, чтобы его удалить.", channel.Title, channel.Name),
	})

	return err
}

func (b *Bot) callbackChannelInfo(btn buttonContext) error {
	var message protobufs.ButtonChanneInfo
	proto.Unmarshal(btn.Data, &message)
	return b.showChannelInfo(btn.Ctx, message.ChannelId, btn.User, btn.UserData)
}

func (b *Bot) callbackBack(btn buttonContext) error {
	message := protobufs.ButtonMenuBack{}
	proto.Unmarshal(btn.Data, &message)

	if message.Msg != nil {
		btn.Data = message.Msg
	}

	return b.btnCallbacks[message.Newmenu](btn)
}

func (b *Bot) showMainPage(ctx context.Context, user *tg.User, userCache *cache.UserData) error {
	_, err := b.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: user.ID},
		ID:          userCache.ActiveMenuID,
		ReplyMarkup: buildInitalMenu(),
		Message:     fmt.Sprintf("Добро пожаловать %s %s", user.FirstName, user.LastName),
	})
	return err
}

func (b *Bot) callbackMainPage(btn buttonContext) error {
	return b.showMainPage(btn.Ctx, btn.User, btn.UserData)
}

func (b *Bot) callbackAddNewKeyWord(btn buttonContext) error {
	var message protobufs.ButtonChanneInfo
	proto.Unmarshal(btn.Data, &message)

	channel, ok := btn.UserData.Channels[message.ChannelId]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при чтении ключевых слов.")
		return nil
	}

	btn.UserData.State = cache.WaitingKeyWord

	rows := []tg.KeyboardButtonRow{CreateBackButton("Отмена", protobufs.MessageID_ChannelInfo, &message)}
	_, err := b.API().MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserData.ActiveMenuID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     fmt.Sprintf("Канал: %s (%s)\nВведите ключевое слово или регулярное выражение.", channel.Title, channel.Name),
	})

	return err
}

func (b *Bot) callbackRemoveKeyWord(btn buttonContext) error {
	var message protobufs.ButtonRemoveKeyWord
	proto.Unmarshal(btn.Data, &message)

	channel, ok := btn.UserData.Channels[message.GroupId]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при удалении.")
		return nil
	}

	err := channel.RemoveKeyword(btn.UserData.TelegramID, message.KeywordId, true)
	if err != nil {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при удалении.")
		return nil
	}

	return b.showChannelInfo(btn.Ctx, channel.TelegramID, btn.User, btn.UserData)
}

func (b *Bot) callbackRemoveChannel(btn buttonContext) error {
	var message protobufs.ButtonChanneInfo
	proto.Unmarshal(btn.Data, &message)

	channel, ok := btn.UserData.Channels[message.ChannelId]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при удалении.")
		return nil
	}

	channelTitle := channel.Title

	if b.channelsCache.RemoveChannelFromUser(btn.UserData.TelegramID, message.ChannelId, true) != nil {
		b.Answer(btn.User).Textf(btn.Ctx, "Ошибка при удалении канал %s.", channelTitle)
		return nil
	}

	b.Answer(btn.User).Textf(btn.Ctx, "Канал %s был удален.", channelTitle)

	if len(btn.UserData.Channels) == 0 {
		b.Answer(btn.User).Text(btn.Ctx, "У вас нет отслеживаемых каналов, вы перемещены в главное меню.")
		return b.callbackMainPage(btn)
	} else {
		return b.callbackMyChannels(btn)
	}
}

func (b *Bot) callbackSpaceButton(btn buttonContext) error {
	texts := []string{
		"Куда вы нажали?", "Больше не нажимайте на кнопку, вы сломаете бота!", "Упс, что-то пошло не так!",
		"Опять?!", "Вы серьезно?", "Я вас предупреждал!", "О, это становится интересно...",
		"Теперь это просто забавно!", "Мне нравится ваша настойчивость!", "Вы действительно упорны!", "Продолжайте, не останавливайтесь!",
		"Нажимайте, как будто от этого зависит жизнь!", "Вы победили, вот ваш приз",
	}

	/* need some funny image or gif */
	textIndex := (btn.UserData.SecretButtonClicks / 3) % len(texts)
	btn.UserData.SecretButtonClicks += 1

	return b.SetAnswerCallback(btn.Ctx, texts[textIndex], btn.Update.QueryID)
}

func (b *Bot) registerQueryCallbacks() {
	b.btnCallbacks[protobufs.MessageID_AddNewChannel] = b.callbackAddNewChannel
	b.btnCallbacks[protobufs.MessageID_MyChannels] = b.callbackMyChannels
	b.btnCallbacks[protobufs.MessageID_AddNewKeyWord] = b.callbackAddNewKeyWord
	b.btnCallbacks[protobufs.MessageID_RemoveKeyWord] = b.callbackRemoveKeyWord
	b.btnCallbacks[protobufs.MessageID_Back] = b.callbackBack
	b.btnCallbacks[protobufs.MessageID_MainPage] = b.callbackMainPage
	b.btnCallbacks[protobufs.MessageID_ChannelInfo] = b.callbackChannelInfo
	b.btnCallbacks[protobufs.MessageID_RemoveChannel] = b.callbackRemoveChannel
	b.btnCallbacks[protobufs.MessageID_Spacer] = b.callbackSpaceButton
}
