package bot

import (
	"context"
	"fmt"
	"tg_reader_bot/internal/app"
	"time"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func (bot *Bot) ParseChannels(ctx context.Context) {
	for {
		app := app.GetContainer()
		if app.Client != nil {
			sender := message.NewSender(app.Client)
			cache := &bot.channelsCache
			cache.Mutex.Lock()
			for _, channelInfo := range cache.Channels {
				channel, err := GetChannelByName(app.Client, sender, ctx, channelInfo.Name)
				if err != nil {
					fmt.Println("GetChannelByName err", err)
					continue
				}

				fmt.Println("Peer LastMsgID", channel)
				msgsClass, err := app.Client.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
					Peer: channel.AsInputPeer(),
				})

				if err != nil {
					fmt.Println("Failed to MessagesGetHistory", err)
					continue
				}

				fmt.Println(msgsClass)

				/*msgs, ok := msgsClass.(*tg.MessagesChannelMessages)
				if !ok {
					fmt.Println("Failed to cast messages")
					continue
				}*/
			}
			cache.Mutex.Unlock()
		}
		time.Sleep(30 * time.Second)
	}

}
