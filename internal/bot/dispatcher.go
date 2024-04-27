package bot

import "tg_reader_bot/internal/events"

func (b *Bot) addCommand(name string, desciption string, callback func(events.Message) error) {
	b.cmdsCallbacks[name] = commandInfo{desciption, callback}
}

func (b *Bot) Dispatch(name string, msg events.Message) error {
	if command, ok := b.cmdsCallbacks[name]; ok {
		return command.callback(msg)
	}
	return nil
}

func (b *Bot) GetCommands() map[string]commandInfo {
	return b.cmdsCallbacks
}
