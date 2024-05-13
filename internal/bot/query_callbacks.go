package bot

import (
	"context"
	"fmt"
	"math"
	"sort"
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/protobufs"

	"github.com/gotd/td/tg"
	"google.golang.org/protobuf/proto"
)

func (b *Bot) callbackAddNewPeer(btn buttonContext) error {
	rows := []tg.KeyboardButtonRow{CreateBackButton("❌ Отмена", protobufs.MessageID_MainPage, nil)}
	btn.UserData.State = cache.WaitingPeerName
	_, err := b.API().MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserData.ActiveMessageID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     "🔗 Введите в чат ссылку/айди имя чата/группы.",
	})
	return err
}

func (b *Bot) showMyPeers(ctx context.Context, userCache *cache.UserData, QueryID int64, page int, sendNewMessage bool) error {
	if len(userCache.Peers) == 0 {
		return b.SetAnswerCallback(ctx, "📄 Список каналов пуст", QueryID)
	}

	maxPages := 0
	pageSize := 5
	maxPages = int(math.Ceil(float64(len(userCache.Peers)) / float64(pageSize)))

	if page < 0 {
		page = 0
	} else if page >= maxPages {
		page = maxPages - 1
	}

	var rows []tg.KeyboardButtonRow
	currentPage := fmt.Sprintf("📄 Страница %d/%d", page+1, maxPages)
	rows = append(rows,
		CreateButtonRow(currentPage, protobufs.MessageID_Spacer, nil),
	)

	startIndex := page * pageSize
	endIndex := startIndex + pageSize

	keys := make([]int64, 0, len(userCache.Peers))
	for k := range userCache.Peers {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	i := 0
	for _, id := range keys {
		peer := userCache.Peers[id]
		if i >= startIndex && i < endIndex {
			rows = append(rows, CreateRowButton(
				peer.Title,
				protobufs.MessageID_PeerInfo,
				&protobufs.ButtonPeerInfo{PeerId: peer.TelegramID, CurrentPage: 0},
			))
		}
		i++
	}

	if maxPages > 0 {
		prevPage := int32(max(page-1, 0))
		nextPage := int32(min(page+1, maxPages-1))
		pagination := tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton("⬅️", protobufs.MessageID_MyPeers, &protobufs.ButtonMyPeers{CurrentPage: prevPage}),
				CreateButton("➡️", protobufs.MessageID_MyPeers, &protobufs.ButtonMyPeers{CurrentPage: nextPage}),
			},
		}
		rows = append(rows, pagination)
	}

	rows = append(rows,
		CreateBackButton("↩️ Назад", protobufs.MessageID_MainPage, nil),
	)

	messageText := "💬 Ваши отслеживаемые каналы, нажмите, чтобы настроить.\n "
	if !sendNewMessage {
		_, err := b.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
			Peer:        &tg.InputPeerUser{UserID: userCache.TelegramID},
			ID:          userCache.ActiveMessageID,
			ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
			Message:     messageText,
		})
		return err
	} else {
		b.DeleteMessage(ctx, userCache.ActiveMessageID)
		_, err := b.Sender.To(&tg.InputPeerUser{UserID: userCache.TelegramID}).Markup(&tg.ReplyInlineMarkup{Rows: rows}).Text(ctx, messageText)
		return err
	}
}

func (b *Bot) callbackMyPeers(btn buttonContext) error {
	var message protobufs.ButtonMyPeers
	proto.Unmarshal(btn.Data, &message)

	return b.showMyPeers(btn.Ctx, btn.UserData, btn.Update.QueryID, int(message.CurrentPage), false)
}

func (b *Bot) showPeerInfo(ctx context.Context, peerID int64, tgUser *tg.User, page int, user *cache.UserData, sendNewMessage bool) error {
	peer, ok := user.Peers[peerID]
	if !ok {
		b.Answer(tgUser).Textf(ctx, "🛑 Ошибка при поиске канала.")
		return nil
	}

	user.ActivePeerID = peerID

	rows := []tg.KeyboardButtonRow{
		CreateRowButton(
			"📝 Новое ключевое слово",
			protobufs.MessageID_AddNewKeyWord,
			&protobufs.ButtonPeerInfo{PeerId: peerID},
		),
	}

	keywords := peer.GetUserKeyWords(user.GetID())

	maxPages := 0
	if keywords != nil {
		pageSize := 5
		maxPages = int(math.Ceil(float64(len(keywords)) / float64(pageSize)))

		if page < 0 {
			page = 0
		} else if page >= maxPages {
			page = maxPages - 1
		}

		currentPage := fmt.Sprintf("📄 Страница %d/%d", page+1, maxPages)
		rows = append(rows,
			CreateButtonRow(currentPage, protobufs.MessageID_Spacer, nil),
		)

		startIndex := page * pageSize
		endIndex := startIndex + pageSize

		keys := make([]int64, 0, len(keywords))
		for k := range keywords {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		i := 0
		for _, id := range keys {
			keyWord := keywords[id]
			if i >= startIndex && i < endIndex {
				rows = append(rows,
					CreateRowButton(
						keyWord,
						protobufs.MessageID_RemoveKeyWord,
						&protobufs.ButtonRemoveKeyWord{
							KeywordId: id,
							PeerInfo:  &protobufs.ButtonPeerInfo{PeerId: peerID, CurrentPage: int32(page)},
						},
					),
				)
			}
			i++
		}
	}

	if maxPages > 0 {
		prevPage := int32(max(page-1, 0))
		nextPage := int32(min(page+1, maxPages-1))
		pagination := tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				CreateButton("⬅️", protobufs.MessageID_PeerInfo, &protobufs.ButtonPeerInfo{PeerId: peerID, CurrentPage: prevPage}),
				CreateButton("➡️", protobufs.MessageID_PeerInfo, &protobufs.ButtonPeerInfo{PeerId: peerID, CurrentPage: nextPage}),
			},
		}
		rows = append(rows, pagination)
	}

	rows = append(rows,
		CreateButtonRow("🗑️ Удалить канал", protobufs.MessageID_RemovePeer, &protobufs.ButtonPeerInfo{PeerId: peerID}),
		CreateBackButton("↩️ Назад", protobufs.MessageID_MyPeers, nil),
		CreateBackButton("⤴️ На главную", protobufs.MessageID_MainPage, nil),
	)

	createStr := peer.CreatedAt.Format("2006-01-02 15:04:05")
	updateStr := peer.UpdatedAt.Format("2006-01-02 15:04:05")

	messageText := fmt.Sprintf(`
💬Канал: %s
📊Ключевых слов: %d
🗓️Добавлен: %s
🗓️Дата последнего обновления: %s
🗑️Нажмите на слово, чтобы его удалить.`, peer.Title, peer.GetUserKeyWordsCount(user.GetID()), createStr, updateStr)

	if !sendNewMessage {
		_, err := b.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
			Peer:        &tg.InputPeerUser{UserID: user.GetID()},
			ID:          user.ActiveMessageID,
			ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
			Message:     messageText,
		})
		return err
	} else {
		b.DeleteMessage(ctx, user.ActiveMessageID)
		_, err := b.Sender.To(&tg.InputPeerUser{UserID: user.GetID()}).Markup(&tg.ReplyInlineMarkup{Rows: rows}).Text(ctx, messageText)
		return err
	}
}

func (b *Bot) callbackPeerInfo(btn buttonContext) error {
	btn.UserData.State = cache.StateNone

	var message protobufs.ButtonPeerInfo
	proto.Unmarshal(btn.Data, &message)
	return b.showPeerInfo(btn.Ctx, message.PeerId, btn.User, int(message.CurrentPage), btn.UserData, false)
}

func (b *Bot) callbackBack(btn buttonContext) error {
	message := protobufs.ButtonMenuBack{}
	proto.Unmarshal(btn.Data, &message)

	if message.Msg != nil {
		btn.Data = message.Msg
	}

	return b.btnCallbacks[message.Newmenu](btn)
}

func (b *Bot) showMainPage(ctx context.Context, user *tg.User, userCache *cache.UserData, sendNewMessage bool) error {
	peersCount := 0
	if userCache != nil {
		peersCount = len(userCache.Peers)
	}

	messageText := fmt.Sprintf(`
✨Добро пожаловать: %s %s ✨
💬Каналов отслеживается: %d 💬
`, user.FirstName, user.LastName, peersCount)

	if !sendNewMessage {
		_, err := b.API().MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
			Peer:        &tg.InputPeerUser{UserID: user.ID},
			ID:          userCache.ActiveMessageID,
			ReplyMarkup: buildInitalMenu(),
			Message:     messageText,
		})
		return err
	} else {
		if userCache != nil {
			b.DeleteMessage(ctx, userCache.ActiveMessageID)
		}
		_, err := b.Sender.To(&tg.InputPeerUser{UserID: user.GetID()}).Markup(buildInitalMenu()).Text(ctx, messageText)
		return err
	}
}

func (b *Bot) callbackMainPage(btn buttonContext) error {
	return b.showMainPage(btn.Ctx, btn.User, btn.UserData, false)
}

func (b *Bot) callbackAddNewKeyWord(btn buttonContext) error {
	var message protobufs.ButtonPeerInfo
	proto.Unmarshal(btn.Data, &message)

	peer, ok := btn.UserData.Peers[message.PeerId]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "🛑 Ошибка при чтении ключевых слов.")
		return nil
	}

	btn.UserData.State = cache.WaitingKeyWord

	rows := []tg.KeyboardButtonRow{CreateBackButton("↩️ Назад", protobufs.MessageID_PeerInfo, &message)}
	_, err := b.API().MessagesEditMessage(btn.Ctx, &tg.MessagesEditMessageRequest{
		Peer:        &tg.InputPeerUser{UserID: btn.Update.UserID},
		ID:          btn.UserData.ActiveMessageID,
		ReplyMarkup: &tg.ReplyInlineMarkup{Rows: rows},
		Message:     fmt.Sprintf("💬Канал: %s \n✍Ввод ключевых слов, каждое новым сообщеним.", peer.Title),
	})

	return err
}

func (b *Bot) callbackRemoveKeyWord(btn buttonContext) error {
	var message protobufs.ButtonRemoveKeyWord
	proto.Unmarshal(btn.Data, &message)

	createNewMenu := false
	peer, ok := btn.UserData.Peers[message.PeerInfo.PeerId]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "🛑 Ошибка при удалении.")
		createNewMenu = true
	}

	err := peer.RemoveKeyword(btn.UserData.GetID(), message.KeywordId)
	if err != nil {
		b.Answer(btn.User).Text(btn.Ctx, "🛑 Ошибка при удалении.")
		createNewMenu = true
	}

	return b.showPeerInfo(btn.Ctx, message.PeerInfo.PeerId, btn.User, int(message.PeerInfo.CurrentPage), btn.UserData, createNewMenu)
}

func (b *Bot) callbackRemovePeer(btn buttonContext) error {
	var message protobufs.ButtonPeerInfo
	proto.Unmarshal(btn.Data, &message)

	peer, ok := btn.UserData.Peers[message.PeerId]
	if !ok {
		b.Answer(btn.User).Text(btn.Ctx, "🛑 Ошибка при удалении.")
	}

	if b.peersCache.RemovePeerFromUser(btn.UserData, peer) != nil {
		b.Answer(btn.User).Textf(btn.Ctx, "🛑 Ошибка при удалении канала %s.", peer.Title)
	}

	b.Answer(btn.User).Textf(btn.Ctx, "✅ %s был удален.", peer.Title)

	if len(btn.UserData.Peers) == 0 {
		b.Answer(btn.User).Text(btn.Ctx, "⚠️ У вас нет отслеживаемых каналов, вы перемещены в главное меню.")
		return b.showMainPage(btn.Ctx, btn.User, btn.UserData, true)
	} else {
		return b.showMyPeers(btn.Ctx, btn.UserData, 0, 0, true)
	}
}

func (b *Bot) callbackSpaceButton(btn buttonContext) error {
	texts := []string{
		"Куда вы нажали?", "Больше не нажимайте на кнопку, вы сломаете бота!", "Упс, что-то пошло не так!",
		"Опять?!", "Вы серьезно?", "Я вас предупреждал!", "О, это становится интересно...",
		"Теперь это просто забавно!", "Мне нравится ваша настойчивость!", "Вы действительно упорны!", "Продолжайте, не останавливайтесь!",
		"Нажимайте, как будто от этого зависит жизнь!", "Вы победили, вот ваш приз",
	}

	/* need funny image or gif */
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
