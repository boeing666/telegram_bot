package bot

import (
	"tg_reader_bot/internal/events"

	"github.com/gotd/td/telegram/message/markup"
)

func (b *Bot) startCommand(msg events.MsgContext) error {
	welcomeText := "Привет, я бот для отслеживания сообщений в чатах и каналах.\n" +
		"Ты можешь добавить необходимый канал, и настроить ключевые слова для него.\n" +
		"Чтобы отслеживать приватные чаты или каналы, добавь меня в них.\n"

	buttons := markup.InlineRow(
		CreateButton(
			"Добавить канал",
			AddNewChannel,
			nil,
		),
		CreateButton(
			"Мои каналы",
			MyChannels,
			nil,
		),
	)

	_, err := b.Answer(msg.PeerUser).Markup(buttons).Text(msg.Ctx, welcomeText)
	return err
}

func (h *Bot) registerCommands() {
	h.addCommand("/start", "Стартовая команда", h.startCommand)
}
