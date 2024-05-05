package bot

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram/message/peer"
	"github.com/gotd/td/tg"
)

func (b *Bot) getChannelByName(ctx context.Context, name string) (*tg.Channel, error) {
	otherPeer, err := b.Sender.Resolve(name).AsInputPeer(ctx)
	if err != nil {
		return nil, err
	}

	switch otherPeer.(type) {
	case *tg.InputPeerChannel:
		inputChannel, ok := peer.ToInputChannel(otherPeer)
		if !ok {
			return nil, fmt.Errorf("cannot cast to ToInputChannel")
		}

		/* maybe use ChannelsGetFullChannel ? */
		channels, err := b.Client.ChannelsGetChannels(ctx, []tg.InputChannelClass{inputChannel})
		if err != nil {
			return nil, err
		}

		chats := channels.GetChats()
		if len(chats) == 0 {
			return nil, fmt.Errorf("getChats return empty slice")
		}

		return chats[0].(*tg.Channel), nil
	default:
		return nil, nil
	}
}
