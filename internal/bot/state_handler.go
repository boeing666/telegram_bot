package bot

import (
	"tg_reader_bot/internal/app"
	"tg_reader_bot/internal/cache"
	"tg_reader_bot/internal/events"
)

func (b *Bot) enterPeerName(msg events.MsgContext) error {
	client := app.GetClient()

	user := msg.UserData
	user.State = cache.StateNone

	b.DeleteMessage(msg.Ctx, msg.Message.ID)

	peer, err := GetChannelByName(client.Client.API(), client.Sender, msg.Ctx, msg.GetText())
	if err != nil {
		b.Answer(msg.PeerUser).Text(msg.Ctx, "Ошибка при выполнении. Попробуйте позже.")
		return err
	}

	if user.HasPeerByID(peer.ID) {
		b.Answer(msg.PeerUser).NoWebpage().Textf(msg.Ctx, "Канал [%s](%s) уже был добавлен.", peer.Title, msg.GetText())
		return err
	}

	err = b.peersCache.AddPeerToUser(msg.UserData, peer)
	if err != nil {
		b.Answer(msg.PeerUser).Text(msg.Ctx, "Ошибка при выполнении. Попробуйте позже.")
		return err
	}

	b.Answer(msg.PeerUser).NoWebpage().Textf(msg.Ctx, "Канал [%s](%s) успешно добавлен.", peer.Title, msg.GetText())
	return b.showPeerInfo(msg.Ctx, peer.ID, msg.PeerUser, user)
}

func (b *Bot) enterKeyWord(msg events.MsgContext) error {
	user := msg.UserData

	b.DeleteMessage(msg.Ctx, msg.Message.ID)

	peer := user.GetActivePeer()
	if peer == nil {
		b.Answer(msg.PeerUser).Textf(msg.Ctx, "Ошибка при получении канала.")
		return nil
	}

	err := peer.AddKeyword(user, msg.GetText())
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) stateHandler(msg events.MsgContext) (bool, error) {
	switch msg.UserData.State {
	case cache.WaitingPeerName:
		return true, b.enterPeerName(msg)
	case cache.WaitingKeyWord:
		return true, b.enterKeyWord(msg)
	}

	return false, nil
}
