package telegram

import (
	"context"
	"tg_reader_bot/internal/commands"
	"tg_reader_bot/internal/config"
	"tg_reader_bot/internal/container"
	"tg_reader_bot/internal/session"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func registerCommandsInBot(ctx context.Context, client *tg.Client) error {
	container := container.GetContainer()
	handlers := container.Handler.Commands

	commands := make([]tg.BotCommand, 0, len(handlers))
	for key, value := range handlers {
		commands = append(commands, tg.BotCommand{
			Command:     key,
			Description: value.Description,
		})
	}

	if _, err := client.BotsSetBotCommands(ctx, &tg.BotsSetBotCommandsRequest{
		Scope:    &tg.BotCommandScopeDefault{},
		Commands: commands,
	}); err != nil {
		return errors.Wrap(err, "register commands")
	}

	return nil
}

func Run(ctx context.Context, config *config.ConfigStructure) error {
	log, _ := zap.NewDevelopment(zap.IncreaseLevel(zapcore.InfoLevel), zap.AddStacktrace(zapcore.FatalLevel))
	defer func() { _ = log.Sync() }()

	dispatcher := tg.NewUpdateDispatcher()

	options := telegram.Options{
		Logger:         log.Named("client"),
		UpdateHandler:  dispatcher,
		SessionStorage: &session.Storage{},
	}

	client := telegram.NewClient(config.AppID, config.AppHash, options)
	sender := message.NewSender(client.API())

	peerManager := peers.Options{
		Logger: log,
	}.Build(client.API())

	gaps := updates.New(updates.Config{
		Handler:      dispatcher,
		AccessHasher: peerManager,
		Logger:       log.Named("gaps"),
	})

	/* listen messages in a channels */
	dispatcher.OnNewChannelMessage(func(ctx context.Context, entities tg.Entities, update *tg.UpdateNewChannelMessage) error {
		return handleChannelMessage(ctx, entities, update, sender)
	})

	/* listen a messages sended to bot, it can be pm or chat */
	dispatcher.OnNewMessage(func(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) error {
		m, ok := update.Message.(*tg.Message)
		if !ok || m.Out {
			return nil
		}

		msg := commands.MessageContext{Ctx: ctx, Entities: entities, Update: update, Sender: sender, Args: nil}

		switch m.PeerID.(type) {
		case *tg.PeerUser: // if msg received in pm
			return handlePrivateMessage(msg)
		case *tg.PeerChat: // if msg received in chat
			return handleGroupChatMessage(msg)
		}

		return nil
	})

	return client.Run(ctx, func(ctx context.Context) error {
		if err := peerManager.Init(ctx); err != nil {
			return err
		}

		/* Check auth status, session maybe is valid */
		status, err := client.Auth().Status(ctx)
		if err != nil {
			return err
		}

		/* Can be already authenticated if we have valid session in session storage. */
		if !status.Authorized {
			if _, err := client.Auth().Bot(ctx, config.APIToken); err != nil {
				return errors.Wrap(err, "auth")
			}
		}

		err = registerCommandsInBot(ctx, client.API())
		if err != nil {
			return err
		}

		u, err := peerManager.Self(ctx)
		if err != nil {
			return err
		}

		_, isBot := u.ToBot()
		return gaps.Run(ctx, client.API(), u.ID(), updates.AuthOptions{
			IsBot: isBot,
		})
	})
}
