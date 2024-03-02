package telegram

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

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

/* Parse all commands here */
func HandlePrivateMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage, sender *message.Sender) error {
	m := update.Message.(*tg.Message)
	_, err := sender.Reply(entities, update).Text(ctx, m.Message)
	return err
}
