package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"tg_reader_bot/internal/app"
	"tg_reader_bot/internal/config"
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

	fmt.Println(config.GetDatabaseQuery())
	db, err := database.Init(config.GetDatabaseQuery())
	if err != nil {
		panic(err)
	}

	app := app.GetContainer()
	app.Init(config, db)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := telegram.Run(ctx, config); err != nil {
		panic(err)
	}
}
