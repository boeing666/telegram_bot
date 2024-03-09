package bot

import (
	"tg_reader_bot/internal/events"

	"github.com/gotd/td/tg"
)

func (b *Bot) startCommand(msg events.Message) error {
	welcomeText := "Привет, я бот для отслеживания сообщений в чатах и каналах\n" +
		"Ты можешь добавить необходимый канал, и настроить ключевые слова для него\n" +
		"Чтобы отслеживать приватные чаты или каналы, добавь меня в них\n"

	buttons := tg.ReplyKeyboardMarkup{
		Rows: []tg.KeyboardButtonRow{{
			Buttons: []tg.KeyboardButtonClass{
				&tg.KeyboardButton{
					Text: "Добавить канал",
				},
				&tg.KeyboardButton{
					Text: "Мои каналы",
				},
			},
		}},
	}

	_, err := b.Answer(msg).Markup(&buttons).Text(msg.Ctx, welcomeText)
	return err
}

func (h *Bot) registerCommands() {
	h.addCommand("/start", "Стартовая команда", h.startCommand)
}
