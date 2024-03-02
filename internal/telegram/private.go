package telegram

import (
	"context"

	"github.com/gotd/td/tg"
	"github.com/k0kubun/pp"
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

func HandlePrivateMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) {
	pp.Println("PrivateMessage")
}
