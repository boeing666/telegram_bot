package telegram

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"github.com/k0kubun/pp"
)

func handleChannelMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewChannelMessage, sender *message.Sender) error {
	pp.Println("ChannelMessage")
	return nil
}
