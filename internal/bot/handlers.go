package bot

import (
	"tg_reader_bot/internal/events"
)

func (b *Bot) handleChannelMessage(msg events.MsgContext) error {

	return nil
}

/* Parse all commands here */
func (b *Bot) handlePrivateMessage(msg events.MsgContext) error {
	if msg.UserCache != nil {
		ok, err := b.stateHandler(msg)
		if ok {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return b.Dispatch(msg.GetMessageText(), msg)
}

func (b *Bot) handleGroupChatMessage(msg events.MsgContext) error {

	return nil
}
