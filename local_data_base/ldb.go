// Задача данного пакета заключается в парсинге базы данных, заноса информации в нее и просчета четных-нечетных недель для последующего заноса в БД.
package local_data_base

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// ////////////////////////////////////////////////////////////////
// Структура содержит массив админов.
type AdminList struct {
	Admins []int64 `json:"admins"`
}

// Функция проверяет является ли переданный пользователь Админом
func IsAdmin(ChatId *int64) bool {
	adminList := AdminList{}
	file, err := os.ReadFile("local_data_base/Jsons/admin_list.json")
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

//////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////

// Структура представляет собой 1 группу, владельца, участников.
type GroupList struct {
	Name      string
	Admin     int64
	Structure []int64
}

// Функция, которая возвращает ссылку на массив структур с группами
func GetGroupsList() (*[]GroupList, error) {
	groupList := &[]GroupList{}
	file, err := os.ReadFile("local_data_base/Jsons/group_list.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, groupList)
	if err != nil {
		return nil, err
	}
	return groupList, nil
}

// Создает группы с указанным именем и добавлей ей админа.
func CreateGroup(groupName string, chatId int64) error {
	groupList := []GroupList{}
	file, err := os.ReadFile("local_data_base/Jsons/group_list.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &groupList)
	if err != nil {
		return err
	}
	for _, v := range groupList {
		if v.Name == groupName {
			return errors.New("такое имя группы уже занято")
		}
	}
	groupList = append(groupList, GroupList{
		Name:  groupName,
		Admin: chatId,
	})
	jsonData, err := json.MarshalIndent(groupList, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile("local_data_base/Jsons/group_list.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Функция удаляет группу
func DeleteGroup(nameGroup string) error {
	groupList, err := GetGroupsList()
	if err != nil {
		return err
	}
	toReturnGropuList := []GroupList{}
	for _, v := range *groupList {
		if v.Name != nameGroup {
			toReturnGropuList = append(toReturnGropuList, v)
		} else {
			err = os.Remove("local_data_base/Jsons/groups/" + nameGroup + ".json")
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
	err = os.WriteFile("local_data_base/Jsons/group_list.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Функция, которая возвращает имя группы по ID из БД.
func IsInside(groupList *[]GroupList, chatId *int64) (string, error) {
	for _, v := range *groupList {
		for _, j := range v.Structure {
			if j == *chatId {
				return v.Name, nil
			}
		}
	}
	return "", errors.New("ChatId: " + fmt.Sprint(*chatId) + " - нет в бд")
}

// Функция выполняет присоединение участника к группе. Предворительно удалив его из предыдущей группы.
func AddUserToGroup(chatId int64, nameGroup string) error {
	groupList, err := GetGroupsList()
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

	for i := range *groupList {
		if (*groupList)[i].Name == nameGroup {
			(*groupList)[i].Structure = append((*groupList)[i].Structure, chatId)
			break
		}

	}
	jsonData, err := json.MarshalIndent(groupList, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile("local_data_base/Jsons/group_list.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

//////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////

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
func GetSchedule(chatId *int64) (*Schedule, error) {
	//Получем список групп
	groupList, err := GetGroupsList()
	if err != nil {
		return &Schedule{}, err
	}
	//Проверяем есть ли там переданный id
	groupName, err := IsInside(groupList, chatId)
	if err != nil {
		return &Schedule{}, errors.New("вы не состоите ни в одной группе")
	}
	//Берем расписание этой группы
	schedules := Schedule{}
	file, err := os.ReadFile("local_data_base/Jsons/groups/" + groupName + ".json")
	if err != nil {
		return &Schedule{}, err
	}
	err = json.Unmarshal(file, &schedules)
	if err != nil {
		return &Schedule{}, err
	}
	return &schedules, nil
}

// Функция создает группу (просчитывает четные-нечетные недели)
func CreateSchedule(groupname string, day, month int, EvenWeekSchedule [][]string, OddWeekSchedule [][]string) error {
	//Получем список групп
	groupList, err := GetGroupsList()
	if err != nil {
		return err
	}
	isInside := false
	for _, v := range *groupList {
		if v.Name == groupname {
			isInside = true
		}
	}
	if !isInside {
		return errors.New("такой группы не в БД")
	}
	even, odd := GetTwoDimensionalArrays(day, month)
	schedule := Schedule{groupname, EvenWeekSchedule, OddWeekSchedule, even, odd}
	// Создание файла для записи
	file, err := os.Create("local_data_base/Jsons/groups/" + groupname + ".json")
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

// Функция находит все даты дней отсноящихся к четной и нечетной неделе отдельно. Аргумент -   число [день][месяц] понедельника нечетной недели.
func GetTwoDimensionalArrays(day, month int) ([][]int, [][]int) {
	monday := time.Date(2024, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	var evenWeeks [][]int
	var oddWeeks [][]int

	for i := 0; i < 55; i++ {
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
