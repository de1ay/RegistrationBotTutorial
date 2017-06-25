package src

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"RegistrationBotTutorial/conf"
)

type TelegramBot struct {
	API                  *tgbotapi.BotAPI        // API телеграмма
	Updates              tgbotapi.UpdatesChannel // Канал обновлений
	ActiveContactRequest []int64                 // ID чатов, от которых мы ожидаем номер
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
		if update.Message != nil && len(update.Message.Text) > 0 {
			// Если сообщение есть и его длина больше 0 -> начинаем обработку
			telegramBot.analyzeUpdate(update)
		}
	}
}

// Начало обработки сообщения
func (telegramBot *TelegramBot) analyzeUpdate(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if telegramBot.findUser(chatID) {
		telegramBot.analyzeUser(chatID)
	} else {
		telegramBot.createUser(User{chatID, ""})
	}
}

func (telegramBot *TelegramBot) findUser(chatID int64) bool {
	find, err := Connection.Find(chatID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать неправильно!")
		telegramBot.API.Send(msg)
	}
	return find
}

func (telegramBot *TelegramBot) createUser(user User) {
	err := Connection.CreateUser(user)
	if err != nil {
		msg := tgbotapi.NewMessage(user.ChatID, "Произошла ошибка! Бот может работать неправильно!")
		telegramBot.API.Send(msg)
	}
}

func (telegramBot *TelegramBot) analyzeUser(chatID int64) {
	user, err := Connection.GetUser(chatID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Произошла ошибка! Бот может работать неправильно!")
		telegramBot.API.Send(msg)
		return
	}
	if len(user.PhoneNumber) > 0 {
		msg := tgbotapi.NewMessage(chatID, "Твой номер: "+ user.PhoneNumber)
		telegramBot.API.Send(msg)
		return
	} else {

	}
}

func (telegramBot *TelegramBot) requestContact(chatID int64) {

}