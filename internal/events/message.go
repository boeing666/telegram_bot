package events

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type Message struct {
	Ctx      context.Context
	Entities tg.Entities
	Update   message.AnswerableMessageUpdate
	Message  *tg.Message
}

func (m *Message) GetMessageText() string {
	return m.Message.Message
}
