package commands

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type commandInfo struct {
	Description string
	callback    func(MessageContext) error
}
type MessageContext struct {
	Ctx      context.Context
	Entities tg.Entities
	Update   *tg.UpdateNewMessage
	Sender   *message.Sender
	Args     []string
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

func (h *Handler) addCommand(name string, desciption string, callback func(MessageContext) error) {
	h.Commands[name] = commandInfo{desciption, callback}
}

func (h *Handler) Handle(name string, msg MessageContext) (bool, error) {
	if command, ok := h.Commands[name]; ok {
		return true, command.callback(msg)
	}
	return false, nil
}
