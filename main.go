package main

import (
	"fmt"
	"os"

	"github.com/mymmrac/telego"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	botUser, err := bot.GetMe()
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("Bot user: %+v\n", botUser)
}
