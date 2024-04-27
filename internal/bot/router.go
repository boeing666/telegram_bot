package bot

import (
	"context"
	"encoding/json"
	"tg_reader_bot/internal/events"

	"github.com/gotd/td/tg"
)

/* listen messages in a channels */
func (b *Bot) onNewChannelMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewChannelMessage) error {
	m, ok := update.Message.(*tg.Message)
	if !ok || m.Out {
		return nil
	}

	peerChannel := m.PeerID.(*tg.PeerChannel)
	tgChannel, ok := entities.Channels[peerChannel.ChannelID]
	if !ok {
		return nil
	}

	msg := events.Message{Ctx: ctx, Entities: entities, Update: update, Message: m, PeerChannel: tgChannel}
	return b.handleChannelMessage(msg)
}

/* listen a messages sended to bot, it can be pm or chat */
func (b *Bot) onNewMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) error {
	m, ok := update.Message.(*tg.Message)
	if !ok || m.Out {
		return nil
	}

	msg := events.Message{Ctx: ctx, Entities: entities, Update: update, Message: m}

	switch m.PeerID.(type) {
	case *tg.PeerUser: // if msg received in pm
		peerUser := m.PeerID.(*tg.PeerUser)
		msg.PeerUser, ok = entities.Users[peerUser.UserID]
		if !ok {
			return nil
		}
		return b.handlePrivateMessage(msg)
	case *tg.PeerChat: // if msg received in chat
		peerChat := m.PeerID.(*tg.PeerChat)
		msg.PeerChat, ok = entities.Chats[peerChat.ChatID]
		if !ok {
			return nil
		}
		return b.handleGroupChatMessage(msg)
	}

	return nil
}

/* called when someone pressed the inline-button */
func (b *Bot) botCallbackQuery(ctx context.Context, entities tg.Entities, update *tg.UpdateBotCallbackQuery) error {
	var queryData QueryHeader
	if err := json.Unmarshal(update.Data, &queryData); err != nil {
		return err
	}

	if queryData.Time < b.startTime {
		_, err := b.Client.MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
			Peer:    &tg.InputPeerUser{UserID: update.UserID},
			ID:      update.MsgID,
			Message: "Сообщение устарело, нажмите /start, чтобы начать работать с ботом.",
		})
		return err
	}

	user, ok := entities.Users[update.UserID]
	if !ok {
		return nil
	}

	userCache := b.getOrCreateUser(update.UserID)
	msg := buttonContext{Ctx: ctx, Entities: entities, Update: update, User: user, UserCache: userCache, Data: queryData.Data}
	if callback, ok := b.btnCallbacks[queryData.Action]; ok {
		return callback(msg)
	}

	return nil
}

func (b *Bot) UpdateHandles(d tg.UpdateDispatcher) {
	d.OnNewChannelMessage(b.onNewChannelMessage)
	d.OnNewMessage(b.onNewMessage)
	d.OnBotCallbackQuery(b.botCallbackQuery)
}
