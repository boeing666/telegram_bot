package telegram

import (
	"context"
	"tg_reader_bot/internal/container"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

/* Parse all commands here */
func handlePrivateMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage, sender *message.Sender) error {
	m := update.Message.(*tg.Message)
	if m.Message[0] == '/' {
		container := container.GetContainer()
		ok, err := container.Handler.Handle(m.Message[1:], ctx, entities, update, sender)
		if err != nil {
			return err
		}
		if !ok {
			_, err := sender.Answer(entities, update).Text(ctx, "Неизвестная команда.")
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
