package bot

import (
	"context"
	"fmt"
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/events"
	"tg_reader_bot/internal/protobufs"
	"time"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type btnCallback func(buttonContext) error

type buttonContext struct {
	Ctx      context.Context
	Entities tg.Entities
	Update   *tg.UpdateBotCallbackQuery
	User     *tg.User
	UserData *cache.UserData
	Data     []byte
}

type commandInfo struct {
	Description string
	callback    func(events.MsgContext) error
}

type Bot struct {
	Client        *tg.Client
	Sender        *message.Sender
	startTime     uint64
	cmdsCallbacks map[string]commandInfo
	btnCallbacks  map[protobufs.MessageID]btnCallback
	channelsCache cache.ChannelsManager
}

func Init(client *tg.Client) *Bot {
	bot := &Bot{
		Client:        client,
		Sender:        message.NewSender(client),
		startTime:     uint64(time.Now().Unix()),
		cmdsCallbacks: make(map[string]commandInfo),
		btnCallbacks:  make(map[protobufs.MessageID]btnCallback),
		channelsCache: cache.ChannelsManager{Channels: make(map[int64]*cache.ChannelInfo), Users: make(map[int64]*cache.UserData)},
	}

	bot.registerCommands()
	bot.registerQueryCallbacks()

	return bot
}

func (b *Bot) LoadUsersChannels() {
	fmt.Println("Loading user channels.")
	b.channelsCache.LoadUsersData()
	fmt.Println("Loading completed.")
}

func (b *Bot) Answer(user *tg.User) *message.RequestBuilder {
	return b.Sender.To(user.AsInputPeer())
}

func (b *Bot) SetAnswerCallback(ctx context.Context, text string, queryID int64) error {
	_, err := b.Client.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: queryID,
		Message: text,
	})
	return err
}

func (b *Bot) DeleteMessage(ctx context.Context, id int) error {
	_, err := b.Client.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
		Revoke: true,
		ID:     []int{id},
	})
	return err
}
