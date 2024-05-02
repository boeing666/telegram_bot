package bot

import "tg_reader_bot/internal/cache"

func (b *Bot) callbackAddNewChannel(btn buttonContext) error {
	btn.UserCache.State = cache.WaitingChannelName
	_, err := b.Sender.To(btn.User.AsInputPeer()).Reply(btn.Update.MsgID).Text(btn.Ctx, "Введите в чат ссылку/айди имя чата/группы.")
	return err
}

func (b *Bot) callbackMyChannels(btn buttonContext) error {

	return nil
}

func (b *Bot) callbackAddNewKeyWord(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackRemoveKeyWord(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackNextChannels(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackPrevChannels(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackNextKeyWords(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackPrevKeyWords(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackBack(btn buttonContext) error {
	return nil
}

func (b *Bot) callbackMainPage(btn buttonContext) error {
	return nil
}

func (b *Bot) registerQueryCallbacks() {
	b.btnCallbacks[AddNewChannel] = b.callbackAddNewChannel
	b.btnCallbacks[MyChannels] = b.callbackMyChannels
	b.btnCallbacks[AddNewKeyWord] = b.callbackAddNewKeyWord
	b.btnCallbacks[RemoveKeyWord] = b.callbackRemoveKeyWord
	b.btnCallbacks[NextChannels] = b.callbackNextChannels
	b.btnCallbacks[PrevChannels] = b.callbackNextKeyWords
	b.btnCallbacks[NextKeyWords] = b.callbackPrevKeyWords
	b.btnCallbacks[Back] = b.callbackBack
	b.btnCallbacks[MainPage] = b.callbackMainPage
}
