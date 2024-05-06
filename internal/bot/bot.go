package bot

import (
	"context"
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/events"
	"tg_reader_bot/internal/protobufs"
	"time"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"github.com/muesli/cache2go"
)

type btnCallback func(buttonContext) error

type buttonContext struct {
	Ctx       context.Context
	Entities  tg.Entities
	Update    *tg.UpdateBotCallbackQuery
	User      *tg.User
	UserCache *cache.UserCache
	Data      []byte
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
	cache         cache2go.CacheTable
}

func Init(client *tg.Client) *Bot {
	bot := &Bot{
		Client:        client,
		Sender:        message.NewSender(client),
		startTime:     uint64(time.Now().Unix()),
		cmdsCallbacks: make(map[string]commandInfo),
		btnCallbacks:  make(map[protobufs.MessageID]btnCallback),
		cache:         *cache2go.Cache("users"),
	}

	bot.registerCommands()
	bot.registerQueryCallbacks()
	return bot
}

func (b *Bot) Answer(user *tg.User) *message.RequestBuilder {
	return b.Sender.To(user.AsInputPeer())
}

func (b *Bot) DeleteMessage(ctx context.Context, id int) error {
	_, err := b.Client.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
		Revoke: true,
		ID:     []int{id},
	})
	return err
}
