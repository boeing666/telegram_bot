package bot

import (
	"context"
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
func (b *Bot) BotCallbackQuery(ctx context.Context, entities tg.Entities, update *tg.UpdateBotCallbackQuery) error {
	/* handle inline callback button here */
	return nil
}

func (b *Bot) UpdateHandles(d tg.UpdateDispatcher) {
	d.OnNewChannelMessage(b.onNewChannelMessage)
	d.OnNewMessage(b.onNewMessage)
	d.OnBotCallbackQuery(b.BotCallbackQuery)
}
