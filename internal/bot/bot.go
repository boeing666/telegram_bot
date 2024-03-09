package bot

import (
	"tg_reader_bot/internal/events"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type commandInfo struct {
	Description string
	callback    func(events.Message) error
}

type Bot struct {
	Client   *tg.Client
	Sender   *message.Sender
	commands map[string]commandInfo
}

func Init(client *tg.Client) *Bot {
	bot := &Bot{
		Client:   client,
		Sender:   message.NewSender(client),
		commands: make(map[string]commandInfo),
	}
	bot.registerCommands()
	return bot
}

func (b *Bot) Answer(msg events.Message) *message.RequestBuilder {
	return b.Sender.Answer(msg.Entities, msg.Update)
}
