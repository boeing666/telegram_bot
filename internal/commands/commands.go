package commands

func startCommand(msg MessageContext) error {
	_, err := msg.Sender.Answer(msg.Entities, msg.Update).Text(msg.Ctx, "Вы ввели start.")
	return err
}

func listenCommand(msg MessageContext) error {
	_, err := msg.Sender.Answer(msg.Entities, msg.Update).Text(msg.Ctx, "Вы ввели listen.")
	return err
}

func unlistenCommand(msg MessageContext) error {
	return nil
}

func printCommand(msg MessageContext) error {
	return nil
}

func (h *Handler) registerCommands() {
	h.addCommand("start", "Стартовая команда помощник", startCommand)
	h.addCommand("listen", "Прослушивать канал", listenCommand)
	h.addCommand("unlisten", "Перестать слушать канал", unlistenCommand)
	h.addCommand("mychannels", "Вывести каналы, которые слушаю", printCommand)
}
