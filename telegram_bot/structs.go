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
		toUserMsg.Text, err = todayOrTomorrow(true, fromUserMsg, toUserMsg)

		if err != nil {
			fmt.Println(err.Error())
		}
	case "завтра":
		toUserMsg.Text, err = todayOrTomorrow(false, fromUserMsg, toUserMsg)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "все":
		err, toUserMsg.Text = all(fromUserMsg, toUserMsg)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "выбор группы":
		//Пункт присоединения вызвывает сообщение, под которым есть всех групп
		err, toUserMsg.Text, toUserMsg.ReplyMarkup = join(fromUserMsg, toUserMsg)
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

// Собирает все расписание для пользователя.
func all(fromUserMsg *tgbotapi.Message, toUserMsg *tgbotapi.MessageConfig) (error, string) {
	toUserMsg.ParseMode = "HTML"
	err, schedule := db.GetSchedule(&fromUserMsg.Chat.ID)
	if err != nil {
		if err.Error() == "Вы не состоите ни в одной группе" {
			return nil, err.Error()
		} else {
			return err, err.Error()
		}
	}
	toReturn := sunEmoji + " <b>Четная неделя:</b>\n"
	for i, v := range schedule.EvenWeekSchedule {
		if i != 0 {
			toReturn += "\n"
		}
		toReturn += dayEmoji + " <b>" + weekday(i) + "</b>:\n"
		for _, k := range v {
			toReturn += checkBoxEmoji + " " + k + "\n"
		}
	}
	toReturn += "\n\n" + moonEmoji + " <b>Нечетная неделя:</b>\n"
	for i, v := range schedule.OddWeekSchedule {
		if i != 0 {
			toReturn += "\n"
		}
		toReturn += dayEmoji + " <b>" + weekday(i) + "</b>:\n"
		for _, k := range v {
			toReturn += checkBoxEmoji + " " + k + "\n"
		}
	}
	return nil, toReturn
}

// Присоединяет пользователя к группе, удаляя его из старой группы.
func join(fromUserMsg *tgbotapi.Message, toUserMsg *tgbotapi.MessageConfig) (error, string, tgbotapi.InlineKeyboardMarkup) {
	err, groupsList := db.GetGroupsList()
	if err != nil {
		return err, "Произошла ошибка (", tgbotapi.NewInlineKeyboardMarkup()
	}
	tempMenu := [][]tgbotapi.InlineKeyboardButton{}
	row := []tgbotapi.InlineKeyboardButton{}
	for i, v := range *groupsList {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(v.Name, v.Name))
		if (i+1)%3 == 0 {
			tempMenu = append(tempMenu, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) <= 2 && len(row) > 0 {
		tempMenu = append(tempMenu, row)
	}
	return nil, "Расписание какой группы Вы хотите видеть?", tgbotapi.InlineKeyboardMarkup{InlineKeyboard: tempMenu}
}
