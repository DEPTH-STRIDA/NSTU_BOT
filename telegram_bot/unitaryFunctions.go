package telegram_bot

import (
	db "NSTU_NN_BOT/local_data_base"
	"errors"
	"fmt"
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
	str = "Я не знаю эту команду.\n\n"
	str += sunEmoji + "Хочешь, добавить такой функционал?\n"
	str += moonEmoji + "Возникли проблемы?\n\n"
	str += "Свяжись со мной: https://t.me/Tichomirov2003"
	return str
}

// Вспомогательная функция, которая возвращает название дня недели по номеру. Нумирация с 0, начиная с понедельника.
func weekdayByName(n time.Weekday) string {
	switch n {
	case time.Monday:
		return "Понедельник"
	case time.Tuesday:
		return "Вторник"
	case time.Wednesday:
		return "Среда"
	case time.Thursday:
		return "Четверг"
	case time.Friday:
		return "Пятница"
	case time.Saturday:
		return "Суббота"
	case time.Sunday:
		return "Воскресенье"
	default:
		return ""
	}
}
func weekdayByIndex(n int) string {
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
func getSchedule(skipDay int, fromUserMsg *tgbotapi.Message, toUserMsg *tgbotapi.MessageConfig) (string, error) {
	//Получае расписание, если пользователь есть в БД
	schedules, err := db.GetSchedule(&fromUserMsg.Chat.ID)
	if err != nil {
		return "К сожалению вас нет в базе данных. Обратитесь в поддержку.", err
	}
	date := time.Now()
	date = date.AddDate(0, 0, skipDay)

	day := date.Day()
	month := int(date.Month())
	dayWeek := int(date.Weekday())
	if dayWeek == 0 {
		dayWeek = 7
	}
	fond := false
	var schedule []string
	toReturn := ""

	toUserMsg.ParseMode = "HTML"
	for _, v := range schedules.EvenWeekDate {
		if v[0] == day && v[1] == month {
			toReturn = sunEmoji + "<b>Четная неделя</b>\n"
			fond = true
			schedule = schedules.EvenWeekSchedule[dayWeek-1]
		}
	}
	if !fond {
		for _, v := range schedules.OddWeekkDate {
			if v[0] == day && v[1] == month {
				toReturn = moonEmoji + "<b>Нечетная неделя</b>\n"
				schedule = schedules.OddWeekSchedule[dayWeek-1]
				fond = true
			}
		}
	}
	if !fond {
		return "Для вашей группы не составлено расписание четных, нечетных недель", errors.New("расписание не найдено в БД. Необходимо составить расписание четных нечетных недель")
	}
	toReturn += clockEmoji + " " + date.Format("02.01") + "\n"
	toReturn += dayEmoji + " <b>" + weekdayByName(date.Weekday()) + "</b>\n\n"
	fmt.Println("today ", int(date.Weekday()))
	isEmpty := true
	for _, v := range schedule {
		if v != "" {
			isEmpty = false
			toReturn += checkBoxEmoji + " " + v + "\n"
		}
	}
	if isEmpty {
		toReturn += checkBoxEmoji + "<b>В этот день ничего нет.</b>"

	}
	return toReturn, nil
}

// Собирает все расписание для пользователя.
func all(fromUserMsg *tgbotapi.Message, toUserMsg *tgbotapi.MessageConfig) (string, error) {
	toUserMsg.ParseMode = "HTML"
	schedule, err := db.GetSchedule(&fromUserMsg.Chat.ID)
	if err != nil {
		return "Не удалось получить расписание вашей группы. Обратитесь в поддержку.", err
	}
	/////////////////////////////////////////////////////////////
	///                   ЧЕТНАЯ НЕДЕЛЯ
	toReturn := sunEmoji + " <b>Четная неделя:</b>\n"
	isEmpty := true
	for i, v := range schedule.EvenWeekSchedule {
		if i != 0 {
			toReturn += "\n"
		}
		isEmpty = true
		toReturn += dayEmoji + " <b>" + weekdayByIndex(i) + "</b>:\n"
		for _, k := range v {
			if k != "" {
				isEmpty = false
				toReturn += checkBoxEmoji + " " + k + "\n"
			}
		}
		if isEmpty {
			toReturn += checkBoxEmoji + "<b>В этот день ничего нет.</b> " + "\n"
		}

	}
	///
	/////////////////////////////////////////////////////////////
	///                 НЕЧЕТНАЯ НЕДЕЛЯ
	toReturn += "\n\n" + moonEmoji + " <b>Нечетная неделя:</b>\n"
	for i, v := range schedule.OddWeekSchedule {
		if i != 0 {
			toReturn += "\n"
		}
		isEmpty = true
		toReturn += dayEmoji + " <b>" + weekdayByIndex(i) + "</b>:\n"
		for _, k := range v {
			if k != "" {
				isEmpty = false
				toReturn += checkBoxEmoji + " " + k + "\n"
			}
		}
		if isEmpty {
			toReturn += checkBoxEmoji + "<b>В этот день ничего нет.</b> " + "\n"
		}
	}
	///
	/////////////////////////////////////////////////////////////
	return toReturn, nil
}

// Присоединяет пользователя к группе, удаляя его из старой группы.
func join() (string, tgbotapi.InlineKeyboardMarkup, error) {
	groupsList, err := db.GetGroupsList()
	if err != nil {
		return "Произошла ошибка (", tgbotapi.NewInlineKeyboardMarkup(), err
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
	return "Расписание какой группы Вы хотите видеть?", tgbotapi.InlineKeyboardMarkup{InlineKeyboard: tempMenu}, nil
}
