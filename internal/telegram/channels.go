package telegram

import (
	"context"

	"github.com/gotd/td/tg"
	"github.com/k0kubun/pp"
)

func HandleChannelMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewChannelMessage) {
	pp.Println("ChannelMessage")
}
