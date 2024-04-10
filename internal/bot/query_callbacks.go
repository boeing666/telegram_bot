package bot

func (b *Bot) callbackAddNewChannel(data queryContext) error {

	return nil
}

func (b *Bot) callbackMyChannels(data queryContext) error {

	return nil
}

func (b *Bot) callbackAddNewKeyWord(data queryContext) error {
	return nil
}

func (b *Bot) callbackRemoveKeyWord(data queryContext) error {
	return nil
}

func (b *Bot) callbackNextChannels(data queryContext) error {
	return nil
}

func (b *Bot) callbackPrevChannels(data queryContext) error {
	return nil
}

func (b *Bot) callbackNextKeyWords(data queryContext) error {
	return nil
}

func (b *Bot) callbackPrevKeyWords(data queryContext) error {
	return nil
}

func (b *Bot) callbackBack(data queryContext) error {
	return nil
}

func (b *Bot) callbackMainPage(data queryContext) error {
	return nil
}

func (b *Bot) registerQueryCallbacks() {
	b.queryCallbacks[AddNewChannel] = b.callbackAddNewChannel
	b.queryCallbacks[MyChannels] = b.callbackMyChannels
	b.queryCallbacks[AddNewKeyWord] = b.callbackAddNewKeyWord
	b.queryCallbacks[RemoveKeyWord] = b.callbackRemoveKeyWord
	b.queryCallbacks[NextChannels] = b.callbackNextChannels
	b.queryCallbacks[PrevChannels] = b.callbackNextKeyWords
	b.queryCallbacks[NextKeyWords] = b.callbackPrevKeyWords
	b.queryCallbacks[Back] = b.callbackBack
	b.queryCallbacks[MainPage] = b.callbackMainPage
}
