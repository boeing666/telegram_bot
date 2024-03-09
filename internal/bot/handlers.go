package bot

import (
	"tg_reader_bot/internal/events"

	"github.com/k0kubun/pp/v3"
)

func (b *Bot) handleChannelMessage(msg events.Message) error {
	pp.Println("ChannelMessage")
	return nil
}

/* Parse all commands here */
func (b *Bot) handlePrivateMessage(msg events.Message) error {
	return b.Dispatch(msg.GetMessageText(), msg)
}

func (b *Bot) handleGroupChatMessage(msg events.Message) error {
	pp.Println("GroupChatMessage")
	return nil
}
