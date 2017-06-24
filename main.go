package main

import (
	"RegistrationBotTutorial/src"
)

var telegramBot src.TelegramBot

func main() {
	telegramBot.Init()
	telegramBot.Start()
}
