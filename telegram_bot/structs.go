// Даннный пакет содержит переменные, функции и структуры, необходимые для работы телеграм-бота.
package telegram_bot

import (
	db "NSTU_NN_BOT/local_data_base"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Функция, которая обрабатывает команды от пользователя.
func command(fromUserMsg *tgbotapi.Message, toUserMsg *tgbotapi.MessageConfig) {
	switch fromUserMsg.Command() {
	case "open":
		//Отправляем пользователю панель с кнопками.
		toUserMsg.Text = "Меню открыто"
		toUserMsg.ReplyMarkup = numericKeyboarAdmin
	case "close":
		//Очищаем панель кнопок.
		toUserMsg.Text = "Меню закрыто"
		toUserMsg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	case "start":
		//Обработка самой начальной команды
		toUserMsg.Text = "Привет! Чтобы увидеть расписание, откройте панель с меню и присоединитесь к своей группе."
	default:
		//Если команда от пользователя не распознается
		toUserMsg.Text = unknownСommand()
	}
}

// Обрабатывает текстовых сообщений от пользователя.
func text(fromUserMsg *tgbotapi.Message, toUserMsg *tgbotapi.MessageConfig) {
	var err error
	//Переводим все буквы в нижний регистр.
	switch strings.ToLower(fromUserMsg.Text) {
	case "сегодня":
		toUserMsg.Text, err = getSchedule(0, fromUserMsg, toUserMsg)

		if err != nil {
			fmt.Println(err.Error())
		}
	case "завтра":
		toUserMsg.Text, err = getSchedule(1, fromUserMsg, toUserMsg)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "все":
		toUserMsg.Text, err = all(fromUserMsg, toUserMsg)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "выбор группы":
		//Пункт присоединения вызвывает сообщение, под которым есть всех групп
		toUserMsg.Text, toUserMsg.ReplyMarkup, err = join()
		if err != nil {
			fmt.Println(err.Error())
		}
	case "меню админа":
		toUserMsg.Text = "Меню закрыто"
		toUserMsg.ReplyMarkup = inlineKeyboard
	default:
		//Если сообщение от пользователя не распознается
		toUserMsg.Text = unknownСommand()
	}

}

// Функция, которая обрабатывает нажатия на кнопки.
func callbackMessage(fromUserMsg *tgbotapi.CallbackConfig, toUserMsg *tgbotapi.MessageConfig) {
	//Других пунктов нет, поэтому тут сразу активируется процесс добавления в группу без какого-либо switch case.
	err := db.AddUserToGroup(toUserMsg.ChatID, fromUserMsg.Text)
	if err != nil {
		fmt.Println(err.Error())
		toUserMsg.Text = "Произошла ошибка. Попробуйте еще раз."
		return
	}
	toUserMsg.Text = "Вы успешно добавлены в группу " + fromUserMsg.Text
}
