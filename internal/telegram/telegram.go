package telegram

import (
	"context"
	"tg_reader_bot/internal/config"
	"tg_reader_bot/internal/session"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Run(ctx context.Context, config *config.ConfigStructure) error {
	log, _ := zap.NewDevelopment(zap.IncreaseLevel(zapcore.InfoLevel), zap.AddStacktrace(zapcore.FatalLevel))
	defer func() { _ = log.Sync() }()

	var (
		dispatcher = tg.NewUpdateDispatcher()
		handler    telegram.UpdateHandler
	)

	options := telegram.Options{
		Logger: log.Named("client"),
		UpdateHandler: telegram.UpdateHandlerFunc(func(ctx context.Context, u tg.UpdatesClass) error {
			return handler.Handle(ctx, u)
		}),
		SessionStorage: &session.Storage{},
	}

	client := telegram.NewClient(config.AppID, config.AppHash, options)

	peerManager := peers.Options{
		Logger: log,
	}.Build(client.API())

	gaps := updates.New(updates.Config{
		Handler:      dispatcher,
		AccessHasher: peerManager,
		Logger:       log.Named("gaps"),
	})

	handler = peerManager.UpdateHook(gaps)

	/* listen a channels messages */
	dispatcher.OnNewChannelMessage(func(ctx context.Context, entities tg.Entities, update *tg.UpdateNewChannelMessage) error {
		HandleChannelMessage(ctx, entities, update)
		return nil
	})

	/* listen a messages sended to bot, it can be pm or group */
	dispatcher.OnNewMessage(func(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) error {
		m, ok := update.Message.(*tg.Message)
		if !ok || m.Out {
			return nil
		}

		switch v := m.PeerID.(type) {
		case *tg.PeerUser: // if msg received in pm
			HandlePrivateMessage(ctx, entities, update)
		case *tg.PeerChat: // if msg received in chat
			HandleGroupChatMessage(ctx, entities, update)
		case *tg.PeerChannel:
			pp.Println("tg.PeerChannel")
		default:
			panic(v)
		}

		return nil
	})

	return client.Run(ctx, func(ctx context.Context) error {
		if _, err := client.Auth().Bot(ctx, config.APIToken); err != nil {
			return errors.Wrap(err, "auth")
		}

		user, err := client.Self(ctx)
		if err != nil {
			return errors.Wrap(err, "call self")
		}

		return gaps.Run(ctx, client.API(), user.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {
				log.Info("Bot Connected")
			},
		})
	})
}
