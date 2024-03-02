package telegram

import (
	"context"

	"github.com/gotd/td/tg"
	"github.com/k0kubun/pp"
)

func HandleGroupChatMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) {
	pp.Println("GroupChatMessage")
}
