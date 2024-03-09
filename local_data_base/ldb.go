// Задача данного пакета заключается в парсинге базы данных, заноса информации в нее и просчета четных-нечетных недель для последующего заноса в БД.
package local_data_base

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// Структура содержит массив админов.
type AdminList struct {
	Admins []int64 `json:"admins"`
}

// Функция проверяет является ли переданный пользователь Админом
func IsAdmin(ChatId *int64) bool {
	adminList := AdminList{}
	file, err := ioutil.ReadFile("local_data_base/Jsons/admin_list.json")
	if err != nil {
		return false
	}
	err = json.Unmarshal(file, &adminList)
	if err != nil {
		return false
	}
	for _, v := range adminList.Admins {
		if v == *ChatId {
			return true
		}
	}
	return false
}

// Структура представляет собой 1 группу, владельца, участников.
type GroupList struct {
	Name      string
	Admin     int64
	Structure []int64
}

// Функция, которая возвращает ссылку на массив структур с группами
func GetGroupsList() (error, *[]GroupList) {
	groupList := &[]GroupList{}
	file, err := ioutil.ReadFile("local_data_base/Jsons/group_list.json")
	if err != nil {
		return err, nil
	}
	err = json.Unmarshal(file, groupList)
	if err != nil {
		return err, nil
	}
	return nil, groupList
}

// Создает группы с указанным именем и добавлей ей админа.
func CreateGroup(groupName string, chatId int64) error {
	groupList := []GroupList{}
	file, err := ioutil.ReadFile("local_data_base/Jsons/group_list.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &groupList)
	if err != nil {
		return err
	}
	for _, v := range groupList {
		if v.Name == groupName {
			return errors.New("Такое имя группы уже занято")
		}
	}
	fmt.Println(groupList)
	groupList = append(groupList, GroupList{
		Name:  groupName,
		Admin: chatId,
	})
	fmt.Println(groupList)
	jsonData, err := json.MarshalIndent(groupList, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("local_data_base/Jsons/group_list.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Функция удаляет группу
func DeleteGroup(nameGroup string) error {
	err, groupList := GetGroupsList()
	if err != nil {
		return err
	}
	toReturnGropuList := []GroupList{}
	for _, v := range *groupList {
		if v.Name != nameGroup {
			toReturnGropuList = append(toReturnGropuList, v)
		} else {
			err = os.Remove("local_data_base/Jsons/" + nameGroup + ".json")
			if err != nil {
				fmt.Println("Ошибка при удалении файла:", err)
				return err
			}

		}
	}
	jsonData, err := json.MarshalIndent(toReturnGropuList, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("local_data_base/Jsons/group_list.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Функция, которая возвращает имя группы по ID из БД.
func IsInside(groupList *[]GroupList, chatId *int64) (error, string) {
	for _, v := range *groupList {
		for _, j := range v.Structure {
			if j == *chatId {
				return nil, v.Name
			}
		}
	}
	return errors.New("ChatId: " + fmt.Sprint(*chatId) + " - нет в бд"), ""
}

// Функция выполняет присоединение участника к группе. Предворительно удалив его из предыдущей группы.
func AddUserToGroup(chatId int64, nameGroup string) error {
	err, groupList := GetGroupsList()
	if err != nil {
		return err
	}

	//Обход групп
	for i := range *groupList {
		// Получаем указатель на текущий элемент
		v := &((*groupList)[i])
		for i, k := range v.Structure {
			if k == chatId {
				if nameGroup == v.Name {
					return nil
				}
				temp := v.Structure[:i]
				temp = append(temp, v.Structure[i+1:]...)
				v.Structure = temp
			}
		}
	}

	for i, _ := range *groupList {
		if (*groupList)[i].Name == nameGroup {
			(*groupList)[i].Structure = append((*groupList)[i].Structure, chatId)
			break
		}

	}
	jsonData, err := json.MarshalIndent(groupList, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("local_data_base/Jsons/group_list.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Шаблон заполнения расписания для группы
type Schedule struct {
	//Имя группы
	GroupName string `json:"group_name"`
	//Двумерные массивы расписания недель четной и нечетной. [день неделеи][занятия]
	EvenWeekSchedule [][]string `json:"Even_week_schedule"`
	OddWeekSchedule  [][]string `json:"Odd_week_schedule"`
	//Двумерный массив хранящий даты дней, относящихся к четным-нечетным неделям. [день][месяц]
	EvenWeekDate [][]int `json:"Even_week_date"`
	OddWeekkDate [][]int `json:"Odd_weekk_date"`
}

// Функция, которая возвращает расписание если он ессть в БД.
func GetSchedule(chatId *int64) (error, *Schedule) {
	//Получем список групп
	err, groupList := GetGroupsList()
	if err != nil {
		return err, &Schedule{}
	}
	//Проверяем есть ли там переданный id
	err, groupName := IsInside(groupList, chatId)
	if err != nil {
		return errors.New("Вы не состоите ни в одной группе"), &Schedule{}
	}
	//Берем расписание этой группы
	schedules := Schedule{}
	file, err := ioutil.ReadFile("local_data_base/Jsons/" + groupName + ".json")
	if err != nil {
		return err, &Schedule{}
	}
	err = json.Unmarshal(file, &schedules)
	if err != nil {
		return err, &Schedule{}
	}
	return nil, &schedules
}
func CreateSchedule(groupname string, day, month int, EvenWeekSchedule [][]string, OddWeekSchedule [][]string) error {
	//Получем список групп
	err, groupList := GetGroupsList()
	if err != nil {
		return err
	}
	isInside := false
	for _, v := range *groupList {
		if v.Name == groupname {
			isInside = true
		}
	}
	if isInside == false {
		return errors.New("Такой группы не в БД.")
	}
	even, odd := GetTwoDimensionalArrays(day, month)
	schedule := Schedule{groupname, EvenWeekSchedule, OddWeekSchedule, even, odd}
	// Создание файла для записи
	file, err := os.Create("local_data_base/Jsons/" + groupname + ".json")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return err
	}
	defer file.Close()

	// Кодирование данных в формат JSON и запись в файл
	encoder := json.NewEncoder(file)
	err = encoder.Encode(schedule)
	if err != nil {
		fmt.Println("Ошибка при кодировании и записи данных:", err)
		return err
	}

	fmt.Println("Данные успешно сохранены в файле" + groupname + ".json")

	return nil
}

// Получить расписание на сегодня
func GetThisOrNexntDay(isToday bool, chatId *int64) (error, []string, string) {
	//Получае расписание, если пользователь есть в БД
	err, schedules := GetSchedule(chatId)
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

// Функция находит все даты дней отсноящихся к четной и нечетной неделе отдельно. Аргумент -   число [день][месяц] понедельника нечетной недели.
func GetTwoDimensionalArrays(day, month int) ([][]int, [][]int) {
	monday := time.Date(2024, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	var evenWeeks [][]int
	var oddWeeks [][]int

	for i := 0; i < 26; i++ {
		for j := 0; j < 7; j++ {
			//Массив Месяц-День
			days := []int{}
			date := monday.AddDate(0, 0, (i*7)+j)
			days = append(days, int(date.Day()))
			days = append(days, int(date.Month()))
			if (i+1)%2 == 0 {
				evenWeeks = append(evenWeeks, days)
			} else {
				oddWeeks = append(oddWeeks, days)
			}
		}

	}
	return evenWeeks, oddWeeks
}
