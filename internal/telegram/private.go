package telegram

import (
	"strings"
	"tg_reader_bot/internal/commands"
	"tg_reader_bot/internal/container"

	"github.com/gotd/td/tg"
)

/* Parse all commands here */
func handlePrivateMessage(msg commands.MessageContext) error {
	m := msg.Update.Message.(*tg.Message)
	msg.Args = strings.Split(m.Message, " ")

	if m.Message[0] == '/' {
		container := container.GetContainer()
		ok, err := container.Handler.Handle(msg.Args[0][1:], msg)
		if err != nil {
			return err
		}
		if !ok {
			_, err := msg.Sender.Answer(msg.Entities, msg.Update).Text(msg.Ctx, "Неизвестная команда.")
			return err
		}
	}

	return nil
}

//api := tg.NewClient(client)
//sender := message.NewSender(api)

/*result, err := peerManager.Resolve(ctx, "@boeingus")
if err != nil {
	return errors.Wrap(err, "resolve")
}

fmt.Println(result.ID())
fmt.Println(result.InputPeer())
fmt.Println(result.InputPeer().TypeName()) */

//_, err := sender.Reply(entities, update).Text(ctx, m.Message)
