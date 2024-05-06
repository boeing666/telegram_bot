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
		_, err := b.Client.MessagesSetBotCallbackAnswer(btn.Ctx, &tg.MessagesSetBotCallbackAnswerRequest{
			QueryID: btn.Update.QueryID,
			Message: "Список каналов пуст",
		})
		return err
	}

	var rows []tg.KeyboardButtonRow
	for _, channel := range btn.UserCache.Channels {
		row := tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					fmt.Sprintf("%s (%d)", channel.Title, len(channel.KeyWords)),
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

	_, err := b.Client.MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserCache.ActiveMenuID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     "Ваши отслеживаемые каналы, нажмите, чтобы настроить.\n ",
	})

	return err
}

func (b *Bot) showChannelInfo(ctx context.Context, channelId int64, User *tg.User, userCache *cache.UserCache) error {
	channel, ok := userCache.Channels[channelId]
	if !ok {
		b.Answer(User).Textf(ctx, "Ошибка при поиске канала.")
		return nil
	}

	userCache.ActiveChannelID = channel.TelegramID

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
		CreateSpaceButtonRow(),
		CreateButtonRow("Удалить канал", protobufs.MessageID_RemoveChannel, &protobufs.ButtonChanneInfo{ChannelId: channel.TelegramID}),
		CreateBackButton("Назад", protobufs.MessageID_MyChannels, nil),
		CreateBackButton("На главную", protobufs.MessageID_MainPage, nil),
	)

	_, err := b.Client.MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: userCache.TelegramID},
		ID:          userCache.ActiveMenuID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     fmt.Sprintf("Ключевые слова для канала %s(%s)\nНажмите на слово, чтобы его удалить.", channel.Title, channel.Name),
	})

	return err
}

func (b *Bot) callbackChannelInfo(btn buttonContext) error {
	var message protobufs.ButtonChanneInfo
	proto.Unmarshal(btn.Data, &message)
	return b.showChannelInfo(btn.Ctx, message.ChannelId, btn.User, btn.UserCache)
}

func (b *Bot) callbackBack(btn buttonContext) error {
	message := protobufs.ButtonMenuBack{}
	proto.Unmarshal(btn.Data, &message)

	if message.Msg != nil {
		btn.Data = message.Msg
	}

	return b.btnCallbacks[message.Newmenu](btn)
}

func (b *Bot) showMainPage(ctx context.Context, user *tg.User, userCache *cache.UserCache) error {
	_, err := b.Client.MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: user.ID},
		ID:          userCache.ActiveMenuID,
		ReplyMarkup: buildInitalMenu(),
		Message:     fmt.Sprintf("Добро пожаловать %s %s", user.FirstName, user.LastName),
	})
	return err
}

func (b *Bot) callbackMainPage(btn buttonContext) error {
	return b.showMainPage(btn.Ctx, btn.User, btn.UserCache)
}

func (b *Bot) callbackAddNewKeyWord(btn buttonContext) error {
	var message protobufs.ButtonChanneInfo
	proto.Unmarshal(btn.Data, &message)

	channel, ok := btn.UserCache.Channels[message.ChannelId]
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
		Message:     fmt.Sprintf("Канал: %s (%s)\nВведите ключевое слово или регулярное выражение.", channel.Title, channel.Name),
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

	return b.showChannelInfo(btn.Ctx, channel.TelegramID, btn.User, btn.UserCache)
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

func (b *Bot) callbackRemoveChannel(btn buttonContext) error {
	var message protobufs.ButtonChanneInfo
	proto.Unmarshal(btn.Data, &message)

	channel, ok := btn.UserCache.Channels[message.ChannelId]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при удалении.")
		return nil
	}

	channelTitle := channel.Title

	if btn.UserCache.RemoveGroup(message.ChannelId) != nil {
		b.Answer(btn.User).Textf(btn.Ctx, "Ошибка при удалении канал %s.", channelTitle)
		return nil
	}

	b.Answer(btn.User).Textf(btn.Ctx, "Канал %s был удален.", channelTitle)

	if len(btn.UserCache.Channels) == 0 {
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
	textIndex := (btn.UserCache.SecretButtonClicks / 3) % len(texts)
	btn.UserCache.SecretButtonClicks += 1

	_, err := b.Client.MessagesSetBotCallbackAnswer(btn.Ctx, &tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: btn.Update.QueryID,
		Message: texts[textIndex],
	})

	return err
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
	b.btnCallbacks[protobufs.MessageID_RemoveChannel] = b.callbackRemoveChannel
	b.btnCallbacks[protobufs.MessageID_Spacer] = b.callbackSpaceButton
}
