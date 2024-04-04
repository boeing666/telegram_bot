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

	return b.handleChannelMessage(events.Message{Ctx: ctx, Entities: entities, Update: update, Message: m})
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
		return b.handlePrivateMessage(msg)
	case *tg.PeerChat: // if msg received in chat
		return b.handleGroupChatMessage(msg)
	}

	return nil
}

/* called when someone pressed the button */
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

	msg := queryContext{Ctx: ctx, Entities: entities, Update: update, Data: queryData.Data}
	if callback, ok := b.queryCallbacks[queryData.Action]; ok {
		return callback(msg)
	}
	return nil
}

func (b *Bot) UpdateHandles(d tg.UpdateDispatcher) {
	d.OnNewChannelMessage(b.onNewChannelMessage)
	d.OnNewMessage(b.onNewMessage)
	d.OnBotCallbackQuery(b.botCallbackQuery)
}

/*
	_, err := b.Client.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: update.QueryID,
		Message: string(update.Data),
	})
*/
