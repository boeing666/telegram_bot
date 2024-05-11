package bot

import (
	"context"
	"fmt"
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/protobufs"

	"github.com/gotd/td/tg"
	"google.golang.org/protobuf/proto"
)

func (b *Bot) callbackAddNewPeer(btn buttonContext) error {
	rows := []tg.KeyboardButtonRow{CreateBackButton("Отмена", protobufs.MessageID_MainPage, nil)}
	btn.UserData.State = cache.WaitingPeerName
	_, err := b.API().MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserData.ActiveMessageID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     "Введите в чат ссылку/айди имя чата/группы.",
	})
	return err
}

func (b *Bot) callbackMyPeers(btn buttonContext) error {
	if len(btn.UserData.Peers) == 0 {
		return b.SetAnswerCallback(btn.Ctx, "Список каналов пуст", btn.Update.QueryID)
	}

	var rows []tg.KeyboardButtonRow
	for peerID, peer := range btn.UserData.Peers {
		row := tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					fmt.Sprintf("(%s) %d", peer.Title, peer.GetUserKeyWordsCount(btn.UserData.GetID())),
					protobufs.MessageID_PeerInfo,
					&protobufs.ButtonPeerInfo{PeerId: peerID},
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
		ID:          btn.UserData.ActiveMessageID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     "Ваши отслеживаемые каналы, нажмите, чтобы настроить.\n ",
	})

	return err
}

func (b *Bot) showPeerInfo(ctx context.Context, peerID int64, User *tg.User, user *cache.UserData) error {
	peer, ok := user.Peers[peerID]
	if !ok {
		b.Answer(User).Textf(ctx, "Ошибка при поиске канала.")
		return nil
	}

	user.ActivePeerID = peerID

	rows := []tg.KeyboardButtonRow{
		{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton(
					"Добавить новое ключевое слово",
					protobufs.MessageID_AddNewKeyWord,
					&protobufs.ButtonPeerInfo{PeerId: peerID},
				),
			},
		},
	}

	keywords := peer.GetUserKeyWords(user.GetID())
	if keywords != nil {
		for id, keyword := range *keywords {
			row := tg.KeyboardButtonRow{
				Buttons: []tg.KeyboardButtonClass{
					CreateButton(
						keyword,
						protobufs.MessageID_RemoveKeyWord,
						&protobufs.ButtonRemoveKeyWord{KeywordId: id, PeerId: peerID},
					),
				},
			}
			rows = append(rows, row)
		}
	}

	rows = append(rows,
		CreateSpaceButtonRow(),
		CreateButtonRow("Удалить канал", protobufs.MessageID_RemovePeer, &protobufs.ButtonPeerInfo{PeerId: peerID}),
		CreateBackButton("Назад", protobufs.MessageID_MyPeers, nil),
		CreateBackButton("На главную", protobufs.MessageID_MainPage, nil),
	)

	_, err := b.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: user.GetID()},
		ID:          user.ActiveMessageID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     fmt.Sprintf("Ключевые слова для канала %s(%s)\nНажмите на слово, чтобы его удалить.", peer.Title, peer.UserName),
	})

	return err
}

func (b *Bot) callbackPeerInfo(btn buttonContext) error {
	var message protobufs.ButtonPeerInfo
	proto.Unmarshal(btn.Data, &message)
	return b.showPeerInfo(btn.Ctx, message.PeerId, btn.User, btn.UserData)
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
		ID:          userCache.ActiveMessageID,
		ReplyMarkup: buildInitalMenu(),
		Message:     fmt.Sprintf("Добро пожаловать %s %s", user.FirstName, user.LastName),
	})
	return err
}

func (b *Bot) callbackMainPage(btn buttonContext) error {
	return b.showMainPage(btn.Ctx, btn.User, btn.UserData)
}

func (b *Bot) callbackAddNewKeyWord(btn buttonContext) error {
	var message protobufs.ButtonPeerInfo
	proto.Unmarshal(btn.Data, &message)

	peer, ok := btn.UserData.Peers[message.PeerId]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при чтении ключевых слов.")
		return nil
	}

	btn.UserData.State = cache.WaitingKeyWord

	rows := []tg.KeyboardButtonRow{CreateBackButton("Отмена", protobufs.MessageID_PeerInfo, &message)}
	_, err := b.API().MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserData.ActiveMessageID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     fmt.Sprintf("Канал: %s (%s)\nМожете писать ключевые слова, каждое слово, новое сообщение\nЧтобы вернутся в меню, нажмите отмена.", peer.Title, peer.UserName),
	})

	return err
}

func (b *Bot) callbackRemoveKeyWord(btn buttonContext) error {
	var message protobufs.ButtonRemoveKeyWord
	proto.Unmarshal(btn.Data, &message)

	peer, ok := btn.UserData.Peers[message.PeerId]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при удалении.")
		return nil
	}

	err := peer.RemoveKeyword(btn.UserData.GetID(), message.KeywordId)
	if err != nil {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при удалении.")
		return nil
	}

	return b.showPeerInfo(btn.Ctx, message.PeerId, btn.User, btn.UserData)
}

func (b *Bot) callbackRemovePeer(btn buttonContext) error {
	var message protobufs.ButtonPeerInfo
	proto.Unmarshal(btn.Data, &message)

	peer, ok := btn.UserData.Peers[message.PeerId]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "Ошибка при удалении.")
		return nil
	}

	if b.peersCache.RemovePeerFromUser(btn.UserData, peer) != nil {
		b.Answer(btn.User).Textf(btn.Ctx, "Ошибка при удалении канал %s.", peer.Title)
		return nil
	}

	b.Answer(btn.User).Textf(btn.Ctx, "Канал %s был удален.", peer.Title)

	if len(btn.UserData.Peers) == 0 {
		b.Answer(btn.User).Text(btn.Ctx, "У вас нет отслеживаемых каналов, вы перемещены в главное меню.")
		return b.callbackMainPage(btn)
	} else {
		return b.callbackMyPeers(btn)
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
	b.btnCallbacks[protobufs.MessageID_AddNewPeer] = b.callbackAddNewPeer
	b.btnCallbacks[protobufs.MessageID_MyPeers] = b.callbackMyPeers
	b.btnCallbacks[protobufs.MessageID_AddNewKeyWord] = b.callbackAddNewKeyWord
	b.btnCallbacks[protobufs.MessageID_RemoveKeyWord] = b.callbackRemoveKeyWord
	b.btnCallbacks[protobufs.MessageID_Back] = b.callbackBack
	b.btnCallbacks[protobufs.MessageID_MainPage] = b.callbackMainPage
	b.btnCallbacks[protobufs.MessageID_PeerInfo] = b.callbackPeerInfo
	b.btnCallbacks[protobufs.MessageID_RemovePeer] = b.callbackRemovePeer
	b.btnCallbacks[protobufs.MessageID_Spacer] = b.callbackSpaceButton
}
