package main

import (
	"context"
	"os"
	"os/signal"
	"tg_reader_bot/internal/commands"
	"tg_reader_bot/internal/config"
	"tg_reader_bot/internal/container"
	"tg_reader_bot/internal/database"
	"tg_reader_bot/internal/session"
	"tg_reader_bot/internal/telegram"
)

func main() {
	err := session.Init()
	if err != nil {
		panic(err)
	}

	config, err := config.Init()
	if err != nil {
		panic(err)
	}

	db, err := database.Init(config.GetDatabaseQuery())
	if err != nil {
		panic(err)
	}

	handler := commands.Init()
	if err != nil {
		panic(err)
	}

	container := container.GetContainer()
	container.Init(config, db, handler)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := telegram.Run(ctx, config); err != nil {
		panic(err)
	}
}
