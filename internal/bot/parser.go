package bot

import (
	"context"
	"fmt"
	"strings"
	"tg_reader_bot/internal/app"
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/events"
	"time"

	"github.com/gotd/td/tg"
)

func (bot *Bot) ParseChannels(ctx context.Context) {
	for {
		app := app.GetContainer()
		if app.Client != nil {
			tgclient := app.Client.Client
			cache := &bot.peersCache
			cache.Mutex.Lock()
			for _, peerInfo := range cache.Peers {
				history, err := tgclient.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
					Peer:  &tg.InputPeerChannel{ChannelID: peerInfo.TelegramID, AccessHash: peerInfo.AccessHash},
					Limit: 10,
					MinID: peerInfo.LastMsgID,
				})

				if err != nil {
					fmt.Println("Failed to MessagesGetHistory", err)
					continue
				}

				modifed, ok := history.AsModified()
				if !ok {
					fmt.Println("Failed to cast to AsModified")
					continue
				}

				messages := modifed.GetMessages()
				if len(messages) == 0 {
					continue
				}

				for _, message := range messages {
					tgmessage, ok := message.(*tg.Message)
					if !ok {
						continue
					}
					bot.FindUsersKeyWords(ctx, tgmessage, peerInfo)
				}

				/* a very crappy lib */
				var id int
				switch v := messages[0].(type) {
				case *tg.MessageEmpty:
					id = v.ID
				case *tg.Message:
					id = v.ID
				case *tg.MessageService:
					id = v.ID
				}

				peerInfo.LastMsgID = id
			}
			cache.Mutex.Unlock()
		}
		time.Sleep(30 * time.Second)
	}
}

func (bot *Bot) ParseIncomingMessage(msg events.MsgContext) {
	cache := &bot.peersCache

	/* a very crappy lib */
	var peerID int64
	switch v := msg.Message.FromID.(type) {
	case *tg.PeerChat:
		peerID = v.ChatID
	case *tg.PeerChannel:
		peerID = v.ChannelID
	}

	peer, ok := cache.Peers[peerID]
	if !ok {
		return
	}

	bot.FindUsersKeyWords(msg.Ctx, msg.Message, peer)
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func (bot *Bot) FindUsersKeyWords(ctx context.Context, message *tg.Message, peerInfo *cache.PeerData) {
	for userID, users := range peerInfo.UsersKeyWords {
		for _, keyword := range users.Keywords {
			if !CaseInsensitiveContains(message.Message, keyword) {
				continue
			}

			bot.Sender.To(&tg.InputPeerUser{UserID: userID}).Textf(ctx, "Найдено ключевое слово: %s\nСообщение: https://t.me/%s/%d", keyword, peerInfo.UserName, message.ID)
			break
		}
	}
}
