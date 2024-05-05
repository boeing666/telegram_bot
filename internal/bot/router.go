package bot

import (
	"context"
	"fmt"
	"tg_reader_bot/internal/events"
	protos "tg_reader_bot/internal/protobufs"

	"github.com/gotd/td/tg"
	"google.golang.org/protobuf/proto"
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

	msg := events.MsgContext{Ctx: ctx, Entities: entities, Update: update, Message: m, PeerChannel: tgChannel}
	return b.handleChannelMessage(msg)
}

/* listen a messages sended to bot, it can be pm or chat */
func (b *Bot) onNewMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) error {
	m, ok := update.Message.(*tg.Message)
	if !ok || m.Out {
		return nil
	}

	msg := events.MsgContext{Ctx: ctx, Entities: entities, Update: update, Message: m}

	switch m.PeerID.(type) {
	case *tg.PeerUser: // if msg received in pm
		peerUser := m.PeerID.(*tg.PeerUser)
		msg.PeerUser, ok = entities.Users[peerUser.UserID]
		msg.UserCache, _ = b.getOrCreateUser(ctx, msg.PeerUser, false)
		if ok {
			return b.handlePrivateMessage(msg)
		}
	case *tg.PeerChat: // if msg received in chat
		peerChat := m.PeerID.(*tg.PeerChat)
		msg.PeerChat, ok = entities.Chats[peerChat.ChatID]
		if ok {
			return b.handleGroupChatMessage(msg)
		}
	}

	return nil
}

/* called when someone pressed the inline-button */
func (b *Bot) botCallbackQuery(ctx context.Context, entities tg.Entities, update *tg.UpdateBotCallbackQuery) error {
	message := protos.MessageHeader{}
	err := proto.Unmarshal(update.Data, &message)
	if err != nil {
		fmt.Println("Unmarshal MessageHeader", err)
		_, err := b.Client.MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
			Peer:    &tg.InputPeerUser{UserID: update.UserID},
			ID:      update.MsgID,
			Message: "Ошибка обработки сообщения, начните заново /start.",
		})
		return err
	}

	if message.Time < b.startTime {
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

	userCache, err := b.getOrCreateUser(ctx, user, true)
	if err != nil {
		b.Answer(user).Textf(ctx, "Ваши данные не загружены, попробуйте позже.")
		return nil
	}

	userCache.SetActiveMenuID(update.MsgID)

	msg := buttonContext{Ctx: ctx, Entities: entities, Update: update, User: user, UserCache: userCache, Data: message.Msg}
	if callback, ok := b.btnCallbacks[message.Msgid]; ok {
		return callback(msg)
	}

	return nil
}

func (b *Bot) UpdateHandles(d tg.UpdateDispatcher) {
	d.OnNewChannelMessage(b.onNewChannelMessage)
	d.OnNewMessage(b.onNewMessage)
	d.OnBotCallbackQuery(b.botCallbackQuery)
}
