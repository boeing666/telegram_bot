package telegram

import (
	"context"
	"tg_reader_bot/internal/bot"
	"tg_reader_bot/internal/config"
	"tg_reader_bot/internal/session"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func registerCommandsInBot(ctx context.Context, bot *bot.Bot) error {
	botCommands := bot.GetCommands()

	commands := make([]tg.BotCommand, 0, len(botCommands))
	for key, value := range botCommands {
		offset := 0
		if key[0] == '/' {
			offset = 1
		}

		commands = append(commands, tg.BotCommand{
			Command:     key[offset:],
			Description: value.Description,
		})
	}

	if _, err := bot.Client.BotsSetBotCommands(ctx, &tg.BotsSetBotCommandsRequest{
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
	updatesHandler := updates.New(updates.Config{
		Handler: dispatcher,
		Logger:  log.Named("gaps"),
	})

	options := telegram.Options{
		Logger:         log.Named("client"),
		UpdateHandler:  updatesHandler,
		SessionStorage: &session.Storage{},
	}

	client := telegram.NewClient(config.AppID, config.AppHash, options)
	api := client.API()

	bot := bot.Init(api)
	bot.UpdateHandles(dispatcher)

	return client.Run(ctx, func(ctx context.Context) error {
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

		user, err := client.Self(ctx)
		if err != nil {
			return errors.Wrap(err, "call self")
		}

		err = registerCommandsInBot(ctx, bot)
		if err != nil {
			return err
		}

		return updatesHandler.Run(ctx, api, user.ID, updates.AuthOptions{
			IsBot:   true,
			OnStart: func(ctx context.Context) { bot.LoadUsersChannels() },
		})
	})
}
