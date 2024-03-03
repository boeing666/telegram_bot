package commands

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func startCommand(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage, sender *message.Sender) error {
	_, err := sender.Answer(entities, update).Text(ctx, "Вы ввели start.")
	return err
}

func listenCommand(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage, sender *message.Sender) error {
	_, err := sender.Answer(entities, update).Text(ctx, "Вы ввели listen.")
	return err
}

func unlistenCommand(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage, sender *message.Sender) error {
	return nil
}

func printCommand(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage, sender *message.Sender) error {
	return nil
}

func (h *Handler) registerCommands() {
	h.addCommand("start", "Стартовая команда помощник", startCommand)
	h.addCommand("listen", "Прослушивать канал", listenCommand)
	h.addCommand("unlisten", "Перестать слушать канал", unlistenCommand)
	h.addCommand("mychannels", "Вывести каналы, которые слушаю", printCommand)
}
