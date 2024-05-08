package bot

import (
	"tg_reader_bot/internal/events"
)

/* Parse all commands here */
func (b *Bot) handlePrivateMessage(msg events.MsgContext) error {
	if msg.UserData != nil {
		ok, err := b.stateHandler(msg)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}
	return b.Dispatch(msg.GetText(), msg)
}

func (b *Bot) handleChannelMessage(msg events.MsgContext) error {

	return nil
}

func (b *Bot) handleGroupChatMessage(msg events.MsgContext) error {

	return nil
}
