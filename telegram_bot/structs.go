// Даннный пакет содержит переменные, функции и структуры, необходимые для работы телеграм-бота.
package telegram_bot

import (
	db "NSTU_NN_BOT/local_data_base"
	"fmt"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Строковые константы с UTF-16 значками.
const (
	clockEmoji    = "\U0001F55A"
	dayEmoji      = "\U0001F4C6"
	checkBoxEmoji = "\u2705"
	sunEmoji      = "\U0001F31E"
	moonEmoji     = "\U0001F31A"
)

// Панель с кнопками.
var numericKeyboarAdmin = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Сегодня"),
		tgbotapi.NewKeyboardButton("Завтра"),
		tgbotapi.NewKeyboardButton("Все"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Выбор группы"),
		tgbotapi.NewKeyboardButton("Меню админа"),
	),
)
var inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Меню админа", os.Getenv("link_web_app")),
	),
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
		toUserMsg.Text = "Я не знаю эту команду.\nХочешь, добавить такой функционал?\n" + sunEmoji + "Свяжись со мной: https://t.me/Tichomirov2003"
	}
}

// Обрабатывает текстовых сообщений от пользователя.
func text(fromUserMsg *tgbotapi.Message, toUserMsg *tgbotapi.MessageConfig) {
	var err error
	//Переводим все буквы в нижний регистр.
	switch strings.ToLower(fromUserMsg.Text) {
	case "сегодня":
		err, toUserMsg.Text = todayOrTomorrow(true, fromUserMsg, toUserMsg)

		if err != nil {
			fmt.Println(err.Error())
		}
	case "завтра":
		err, toUserMsg.Text = todayOrTomorrow(false, fromUserMsg, toUserMsg)
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
		toUserMsg.Text = "Я не знаю эту команду.\nХочешь, добавить такой функционал?\n" + sunEmoji + "Свяжись со мной: https://t.me/Tichomirov2003"
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

// Функция, которая возвращает расписание на этот день
func todayOrTomorrow(isToday bool, fromUserMsg *tgbotapi.Message, toUserMsg *tgbotapi.MessageConfig) (error, string) {
	//Режим чтения строковых перменных. Позволяет использовать метки HTML.
	toUserMsg.ParseMode = "HTML"

	err, schedule, evenOrOdd := db.GetThisOrNexntDay(isToday, &fromUserMsg.Chat.ID)
	if err != nil {
		if err.Error() == "Вы не состоите ни в одной группе" {
			return nil, err.Error()
		} else {
			return err, err.Error()
		}

	}
	date := time.Now()
	if isToday == false {
		date = date.AddDate(0, 0, 1)
	}
	toReturn := moonEmoji + "<b>Нечетная неделя</b>\n"
	if evenOrOdd == "even" {
		toReturn = sunEmoji + "<b>Четная неделя</b>\n"
	}
	toReturn += clockEmoji + " " + time.Now().Format("02.01 - 15:04") + "\n"
	toReturn += dayEmoji + " <b>" + weekday(int(date.Weekday())-1) + "</b>\n\n"
	for _, v := range schedule {
		toReturn += checkBoxEmoji + " " + v + "\n"
	}
	return nil, toReturn
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

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////             переписать на iota const         //////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Вспомогательная функция, которая возвращает название дня недели по номеру. Нумирация с 0, начиная с понедельника.
func weekday(n int) string {
	switch n {
	case 0:
		return "Понедельник"
	case 1:
		return "Вторник"
	case 2:
		return "Среда"
	case 3:
		return "Четверг"
	case 4:
		return "Пятница"
	case 5:
		return "Суббота"
	case 6:
		return "Воскресенье"
	default:
		return ""
	}
}
