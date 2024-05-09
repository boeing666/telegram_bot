package main

import (
	"context"
	"tg_reader_bot/internal/app"
	"tg_reader_bot/internal/bot"
	"tg_reader_bot/internal/client"
	"tg_reader_bot/internal/config"
	"tg_reader_bot/internal/database"
)

func main() {
	config, err := config.Init()
	if err != nil {
		panic(err)
	}

	db, err := database.Init(config.GetDatabaseQuery())
	if err != nil {
		panic(err)
	}

	app := app.GetContainer()
	app.Init(config, db)

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	go bot.Run(context)
	go client.Run(context)
}
