package bot

import (
	"context"
	"tg_reader_bot/internal/events"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type queryCallback func(queryContext) error

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

type queryContext struct {
	Ctx      context.Context
	Entities tg.Entities
	Update   *tg.UpdateBotCallbackQuery
	Data     string
}

type CallbackData struct {
	CreatedTime int64
	ActionType  uint32
	Data        []byte
}

type commandInfo struct {
	Description string
	callback    func(events.Message) error
}

type Bot struct {
	Client         *tg.Client
	Sender         *message.Sender
	startTime      uint64
	commands       map[string]commandInfo
	queryCallbacks map[uint32]queryCallback
	cache          bigcache.BigCache
}

func Init(client *tg.Client) *Bot {
	bot := &Bot{
		Client:         client,
		Sender:         message.NewSender(client),
		startTime:      uint64(time.Now().Unix()),
		commands:       make(map[string]commandInfo),
		queryCallbacks: make(map[uint32]queryCallback),
	}
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	bot.cache = *cache
	bot.registerCommands()
	bot.registerQueryCallbacks()
	return bot
}

func (b *Bot) Answer(msg events.Message) *message.RequestBuilder {
	return b.Sender.Answer(msg.Entities, msg.Update)
}
