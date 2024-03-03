package commands

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type commandInfo struct {
	Description string
	callback    func(context.Context, tg.Entities, *tg.UpdateNewMessage, *message.Sender) error
}

type Handler struct {
	Commands map[string]commandInfo
}

func Init() *Handler {
	var handler Handler
	handler.Commands = make(map[string]commandInfo)
	handler.registerCommands()
	return &handler
}

func (h *Handler) addCommand(name string, desciption string, callback func(context.Context, tg.Entities, *tg.UpdateNewMessage, *message.Sender) error) {
	h.Commands[name] = commandInfo{desciption, callback}
}

func (h *Handler) Handle(name string, ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage, sender *message.Sender) (bool, error) {
	if command, ok := h.Commands[name]; ok {
		return true, command.callback(ctx, entities, update, sender)
	}
	return false, nil
}
