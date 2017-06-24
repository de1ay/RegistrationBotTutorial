package src

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"RegistrationBotTutorial/conf"
	"log"
)

type TelegramBot struct {
	API *tgbotapi.BotAPI  // API телеграмма
	Updates tgbotapi.UpdatesChannel  // Канал обновлений
	ActiveContactRequest []int64  // ID чатов, от которых мы ожидаем номер
}

// Инициализация бота
func (telegramBot *TelegramBot) Init() {
	botAPI, err := tgbotapi.NewBotAPI(conf.TELEGRAM_BOT_API_KEY)  // Инициализация API
	if err != nil {
		log.Fatal(err)
	}
	telegramBot.API = botAPI
	botUpdate := tgbotapi.NewUpdate(conf.TELEGRAM_BOT_UPDATE_OFFSET)  // Инициализация канала обновлений
	botUpdate.Timeout = conf.TELEGRAM_BOT_UPDATE_TIMEOUT
	botUpdates, err := telegramBot.API.GetUpdatesChan(botUpdate)
	if err != nil {
		log.Fatal(err)
	}
	telegramBot.Updates = botUpdates
}

// Основной цикл бота
func (telegramBot *TelegramBot) Start() {
	for update := range telegramBot.Updates {
		if update.Message != nil && len(update.Message.Text) > 0{  // Если сообщение есть и его длина больше 0 -> начинаем обработку
			telegramBot.analyzeUpdate(update)
		}
	}
}

// Начало обработки сообщения
func (telegramBot *TelegramBot) analyzeUpdate(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	telegramBot.API.Send(msg)
}