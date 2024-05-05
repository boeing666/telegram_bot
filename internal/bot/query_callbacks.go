package bot

import (
	"fmt"
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/serializer"

	"github.com/gotd/td/tg"
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
					fmt.Sprintf("%s (%d)", channel.Name, len(channel.KeyWords)),
					ChannelInfo,
					buttonChannelInfo{Name: channel.Name},
				),
			},
		}
		rows = append(rows, row)
	}

	rows = append(rows, tg.KeyboardButtonRow{
		Buttons: []tg.KeyboardButtonClass{
			CreateButton(
				"Назад",
				Back,
				buttonMenuBack{BackMenu: 1},
			),
		},
	})

	_, err := b.Client.MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.Update.MsgID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     "Ваши отслеживаемые каналы, нажмите, чтобы настроить.",
	})

	return err
}

func (b *Bot) callbackChannelInfo(btn buttonContext) error {
	message := buttonChannelInfo{}
	err := serializer.DecodeMessage(btn.Data, &message)

	if err != nil {
		b.Answer(btn.User).Textf(btn.Ctx, "Ошибка при чтении ключевых слов.")
		return err
	}

	channel, ok := btn.UserCache.Channels[message.Name]
	if !ok {
		b.Answer(btn.User).Textf(btn.Ctx, "Ошибка при чтении ключевых слов.")
		return nil
	}

	rows := []tg.KeyboardButtonRow{
		{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					"Добавить новое",
					AddNewKeyWord,
					buttonChannelInfo{Name: channel.Name},
				),
			},
		},
	}

	for id, keyword := range channel.KeyWords {
		row := tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					keyword,
					RemoveKeyWord,
					buttonRemoveKeyWord{ID: id},
				),
			},
		}
		rows = append(rows, row)
	}

	_, err = b.Client.MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.Update.MsgID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     fmt.Sprintf("Ключевые слова для канала %s, нажмите на слово, чтобы удалить.", channel.Name),
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

func (b *Bot) callbackBack(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackMainPage(btn buttonContext) error {
	return nil
}

func (b *Bot) registerQueryCallbacks() {
	b.btnCallbacks[AddNewChannel] = b.callbackAddNewChannel
	b.btnCallbacks[MyChannels] = b.callbackMyChannels
	b.btnCallbacks[AddNewKeyWord] = b.callbackAddNewKeyWord
	b.btnCallbacks[RemoveKeyWord] = b.callbackRemoveKeyWord
	b.btnCallbacks[NextChannels] = b.callbackNextChannels
	b.btnCallbacks[PrevChannels] = b.callbackNextKeyWords
	b.btnCallbacks[NextKeyWords] = b.callbackPrevKeyWords
	b.btnCallbacks[Back] = b.callbackBack
	b.btnCallbacks[MainPage] = b.callbackMainPage
	b.btnCallbacks[ChannelInfo] = b.callbackChannelInfo
}
