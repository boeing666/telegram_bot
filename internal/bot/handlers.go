package bot

import (
	"tg_reader_bot/internal/events"
)

func (b *Bot) handleChannelMessage(msg events.Message) error {

	return nil
}

/* Parse all commands here */
func (b *Bot) handlePrivateMessage(msg events.Message) error {

	return b.Dispatch(msg.GetMessageText(), msg)
}

func (b *Bot) handleGroupChatMessage(msg events.Message) error {

	return nil
}
