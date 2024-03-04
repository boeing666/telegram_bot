package telegram

import (
	"tg_reader_bot/internal/commands"

	"github.com/k0kubun/pp"
)

func handleGroupChatMessage(msg commands.MessageContext) error {
	pp.Println("GroupChatMessage")
	return nil
}
