package events

import (
	"context"
	"tg_reader_bot/internal/cache"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type MsgContext struct {
	Ctx         context.Context
	Entities    tg.Entities
	Update      message.AnswerableMessageUpdate
	Message     *tg.Message
	PeerUser    *tg.User
	PeerChat    *tg.Chat
	PeerChannel *tg.Channel
	UserCache   *cache.UserCache
}

func (m *MsgContext) GetMessageText() string {
	return m.Message.Message
}
