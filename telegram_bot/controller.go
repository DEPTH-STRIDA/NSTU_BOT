// Даннный пакет содержит телеграм-бота.
package telegram_bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Запуск бота
// Данный go - файл всего-лишь выполняет роль маршрутизатора, который распреляет как следует обрабатывать сообщения. Текст, команда, нетекстовые(нажатие на кнопки)
func CreateTgBot(Token string) {
	if Token == "" {
		fmt.Println("Пустой токен бота. Скорее всего не найдена переменная среды. controller.go #15")
		return
	}
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		//log.Panic(err)
		fmt.Println(err.Error())
	}

	//Создание канала, куда приходят сообщения для обработки
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	//Обработка команд от бота
	for update := range updates {
		//Обработка текстовых сообщений: текст, команды.
		if update.Message != nil {
			//Создание структуры ответного сообщения
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			//Вывод в консоль id отправляющего и его команду.
			fmt.Println(update.Message.Chat.ID, ": ", update.Message.Text)

			///////////////////////////////////////////////////////////////
			///                  ОБРАБОТКА КОМАНДЫ                      ///
			///////////////////////////////////////////////////////////////
			if update.Message.IsCommand() {
				command(update.Message, &msg)
			} else {
				///////////////////////////////////////////////////////////////
				///                  ОБРАБОТКА ТЕКСТА                       ///
				///////////////////////////////////////////////////////////////
				text(update.Message, &msg)
			}
			//Отправляем сообщение, если нет ошибок при отправке.
			if _, err := bot.Send(msg); err != nil {
				fmt.Println(err.Error())
			}
			///////////////////////////////////////////////////////////////
			///                   ОБРАБОТКА КНОПОК                      ///
			///////////////////////////////////////////////////////////////
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				fmt.Println(err.Error())
			}
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
			//Обработка команды.
			callbackMessage(&callback, &msg)
			//Отправка сообщения
			if _, err := bot.Send(msg); err != nil {
				fmt.Println(err.Error())
			}
		}

	}
}
