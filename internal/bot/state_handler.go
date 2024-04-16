package bot

import "tg_reader_bot/internal/events"

func (b *Bot) stateHandler(msg events.Message) (bool, error) {
	user := b.getOrCreateUser(msg.PeerUser.ID)
	if user.state == StateNone {
		return false, nil
	}
	switch user.state {
	case WaitingChannelName:

	}
	return true, nil
}
