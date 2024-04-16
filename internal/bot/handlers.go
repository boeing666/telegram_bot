package bot

import (
	"tg_reader_bot/internal/events"
)

func (b *Bot) handleChannelMessage(msg events.Message) error {

	return nil
}

/* Parse all commands here */
func (b *Bot) handlePrivateMessage(msg events.Message) error {
	ok, err := b.stateHandler(msg)
	if ok {
		return nil
	}

	if err != nil {
		return err
	}

	return b.Dispatch(msg.GetMessageText(), msg)
}

func (b *Bot) handleGroupChatMessage(msg events.Message) error {

	return nil
}
