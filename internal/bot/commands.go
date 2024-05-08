package bot

import (
	"tg_reader_bot/internal/events"
)

func (b *Bot) startCommand(msg events.MsgContext) error {
	welcomeText := "Привет, я бот для отслеживания сообщений в чатах и каналах.\n" +
		"Ты можешь добавить необходимый канал, и настроить ключевые слова для него.\n" +
		"Чтобы отслеживать приватные чаты или каналы, добавь меня в них.\n"

	_, err := b.Sender.To(msg.PeerUser.AsInputPeer()).Markup(buildInitalMenu()).Text(msg.Ctx, welcomeText)
	return err
}

func (h *Bot) registerCommands() {
	h.addCommand("/start", "Стартовая команда", h.startCommand)
}
