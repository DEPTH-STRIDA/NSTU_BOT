package telegram_bot

import (
	db "NSTU_NN_BOT/local_data_base"
	"errors"
	"os"
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

// Обработка неизвестных команд.
func unknownСommand() (str string) {
	str = "Я не знаю эту команду.\n"
	str += sunEmoji + "Хочешь, добавить такой функционал?\n"
	str += moonEmoji + "Возникли проблемы?\n"
	str += "Свяжись со мной: https://t.me/Tichomirov2003"
	return str
}

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

// Функция, которая возвращает расписание на этот день
func todayOrTomorrow(isToday bool, fromUserMsg *tgbotapi.Message, toUserMsg *tgbotapi.MessageConfig) (string, error) {
	//Режим чтения строковых перменных. Позволяет использовать метки HTML.
	toUserMsg.ParseMode = "HTML"

	err, schedule, evenOrOdd := db.GetThisOrNexntDay(isToday, &fromUserMsg.Chat.ID)
	if err != nil {
		if err.Error() == "Вы не состоите ни в одной группе" {
			return err.Error(), nil
		} else {
			return err.Error(), err
		}

	}
	date := time.Now()
	if !isToday {
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
	return toReturn, nil
}

// Получить расписание на сегодня
func GetThisOrNexntDay(isToday bool, chatId *int64) (error, []string, string) {
	//Получае расписание, если пользователь есть в БД
	err, schedules := db.GetSchedule(chatId)
	if err != nil {
		return err, nil, ""
	}
	date := time.Now()
	if isToday == false {
		date = date.AddDate(0, 0, 1)
	}

	day := date.Day()
	month := int(date.Month())
	dayWeek := int(date.Weekday())
	if dayWeek == 0 {
		dayWeek = 7
	}

	for _, v := range schedules.EvenWeekDate {
		if v[0] == day && v[1] == month {
			return nil, schedules.EvenWeekSchedule[dayWeek-1], "even"
		}
	}
	for _, v := range schedules.OddWeekkDate {
		if v[0] == day && v[1] == month {
			return nil, schedules.OddWeekSchedule[dayWeek-1], "odd"
		}
	}
	return errors.New("Вы не состоите ни в одной группе"), nil, ""
}
