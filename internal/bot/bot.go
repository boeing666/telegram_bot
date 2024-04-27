package bot

import (
	"context"
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/events"
	"time"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"github.com/muesli/cache2go"
)

type btnCallback func(buttonContext) error

const (
	AddNewChannel = iota
	MyChannels
	AddNewKeyWord
	RemoveKeyWord
	NextChannels
	PrevChannels
	NextKeyWords
	PrevKeyWords
	Back
	MainPage
)

type QueryHeader struct {
	Time   uint64
	Action uint32
	Data   string
}

type buttonContext struct {
	Ctx       context.Context
	Entities  tg.Entities
	Update    *tg.UpdateBotCallbackQuery
	User      *tg.User
	UserCache *cache.UserCache
	Data      string
}

type buttonData struct {
	CreatedTime int64
	ActionType  uint32
	Data        []byte
}

type commandInfo struct {
	Description string
	callback    func(events.Message) error
}

type Bot struct {
	Client        *tg.Client
	Sender        *message.Sender
	startTime     uint64
	cmdsCallbacks map[string]commandInfo
	btnCallbacks  map[uint32]btnCallback
	cache         cache2go.CacheTable
}

func Init(client *tg.Client) *Bot {
	bot := &Bot{
		Client:        client,
		Sender:        message.NewSender(client),
		startTime:     uint64(time.Now().Unix()),
		cmdsCallbacks: make(map[string]commandInfo),
		btnCallbacks:  make(map[uint32]btnCallback),
		cache:         *cache2go.Cache("users"),
	}

	bot.registerCommands()
	bot.registerQueryCallbacks()
	return bot
}

func (b *Bot) Answer(msg events.Message) *message.RequestBuilder {
	return b.Sender.Answer(msg.Entities, msg.Update)
}
