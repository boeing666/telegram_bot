package telegram

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"github.com/k0kubun/pp"
)

func HandleGroupChatMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage, sender *message.Sender) error {
	pp.Println("GroupChatMessage")
	return nil
}
