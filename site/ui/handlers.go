package site

//Задача данного файла в обработки соответсвующих URL путей.
import (
	"NSTU_NN_BOT/local_data_base"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"
)

////////////////////////////////////////////////////////
///                    УНИТАРНОЕ                     ///
////////////////////////////////////////////////////////

type Data struct {
	TextSlice []string
}

// Недопустимые слова для названий групп
var unacceptable_words = [...]string{
	"сегодня", "завтра", "все", "выбор группы", "меню админа", "/open", "/close", "/start",
	"Сегодня", "Завтра", "Все", "Выбор группы", "Меню админа", "/open", "/close", "/start",
}

func createError(str string) error {
	pc, _, _, _ := runtime.Caller(1)
	callerFunction := runtime.FuncForPC(pc).Name()
	_, _, line, _ := runtime.Caller(1)
	return errors.New(str + " callerFunction: " + callerFunction + " line: " + fmt.Sprint(line))
}

////////////////////////////////////////////////////////
///                    СТРАНИЦЫ                      ///
////////////////////////////////////////////////////////

// обработка и верификация telegram пользователя
func (app *Application) validating(w http.ResponseWriter, r *http.Request) {
	app.render(w, "valid.page.tmpl")
}

// Главная страница - список групп
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	// Извлечение параметра "userID" из строки запроса
	id := r.URL.Query().Get("id")
	//Парсин strint в int64
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", createError("ошибка конвертации id, возможно пользователь зашел не из telegram."), debug.Stack())
		app.ErrorLog.Output(2, trace)
		return
	}
	//Если пользователь не является админом, высылаем ему заглушку
	if !local_data_base.IsAdmin(&intId) {
		app.render(w, "noAdmin.page.tmpl")
		return
	}
	//Получаем список групп
	groupList, err := local_data_base.GetGroupsList()
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", createError("Ошибка получения списка групп: "+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		app.serverError(w, createError("Ошибка получения списка групп: "+err.Error()))
	}
	//Формируем массив из групп, которые принадлежат только пользователю
	textSlice := []string{}
	for _, v := range *groupList {
		if v.Admin == intId {
			textSlice = append(textSlice, v.Name)
		}

	}
	//Формируем структур, чтобы передать ее в шаблон
	data := Data{textSlice}
	ts, ok := app.TemplateCache["home.page.tmpl"]
	if !ok {
		trace := fmt.Sprintf("%s\n%s", createError("шаблон home.page.tmpl не существует"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		app.serverError(w, createError("шаблон home.page.tmpl не существует"))
		return
	}
	//Выполнняем шаблон
	err = ts.Execute(w, &data)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", createError("Не удалось выполнить выгрузку шаблона с ошибкой: "+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		app.serverError(w, err)
	}
	app.InfoLog.Output(2, "На главную страницу успешно зашел пользователь. "+string(id))
	app.render(w, "home.page.tmpl")
}

// Страница с добавлением новой группы.
func (app *Application) new_group(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Output(2, "На страницу создания группы успешно зашел пользователь. ")
	app.render(w, "new_group.page.tmpl")
}

////////////////////////////////////////////////////////
///                     ЗАПРОСЫ                      ///
////////////////////////////////////////////////////////

// Обработка аякс запросов. Проверяет данные и возвращает успешный код.
func (app *Application) validate(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodPost {
		trace := fmt.Sprintf("%s\n%s", createError("invalid request method"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, createError("Invalid request method").Error(), http.StatusMethodNotAllowed)
		return
	}
	//Получаем тело запроса
	err := r.ParseForm()
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", createError("error parsing form data (body)"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, createError("Error parsing form data").Error(), http.StatusBadRequest)
		return
	}

	/////////////////////////////////////////////////////////////////////
	///                  ДОБАВИТЬ ОБРАТОКУ ТОКЕНА                     ///
	/////////////////////////////////////////////////////////////////////
	//Парсим ключ из body
	token := r.Form.Get("token")

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)
	app.InfoLog.Output(2, "На сайт успешно зашел пользователь. ")
	fmt.Fprintf(w, "Token received successfully: %s", token)
}

// Обработка запросов на корректность имени группы.
func (app *Application) checkGroupName(w http.ResponseWriter, r *http.Request) {
	//Провекра типа запроса
	if r.Method != http.MethodPost {
		trace := fmt.Sprintf("%s\n%s", createError("Получен неправильный тип запроса"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	//Парсинг тела запроса
	err := r.ParseForm()
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", createError("Ошибка при парсинге тела запроса"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, "Имя использует запрещенные символы.", http.StatusBadRequest)
		return
	}
	//Парсинг переменых из тела запроса
	groupName := r.Form.Get("groupName")
	id := r.Form.Get("id")
	//Получение списка групп
	groupList, err := local_data_base.GetGroupsList()
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", createError("Ошибка при получении списка групп:"+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		w.WriteHeader(http.StatusBadGateway)
		http.Error(w, "Неизвестная ошибка. Подождите или сообщите в тех. поддержку.", http.StatusBadRequest)
		return
	}
	//Проверка предложенного имени на корректность (наличие недопустимых слов)
	for _, v := range unacceptable_words {
		if v == groupName {
			trace := fmt.Sprintf("%s\n%s", createError("Имя группы: "+groupName+" - недопустимо"), debug.Stack())
			app.ErrorLog.Output(2, trace)
			w.WriteHeader(http.StatusBadGateway)
			http.Error(w, "Имя группы \""+groupName+"\" недопустимо\n", http.StatusBadRequest)
			return
		}
	}
	//Проверка предложенного имени на корректность (занятость имени группы)
	for _, v := range *groupList {
		if v.Name == groupName {
			trace := fmt.Sprintf("%s\n%s", createError("Имя группы: "+groupName+" - занято"), debug.Stack())
			app.ErrorLog.Output(2, trace)
			w.WriteHeader(http.StatusBadGateway)
			http.Error(w, "Имя группы \""+groupName+"\" занято", http.StatusBadRequest)
			return
		}
	}
	// Конвертация строки в int64
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", createError("Ошибка при конвертации id: "+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Ошибка при конвертации id")
		return
	}
	//Создание предложенной группы в БД
	err = local_data_base.CreateGroup(groupName, intId)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", createError("Ошибка при добавлении группы: "+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Ошибка при добавлении группы")
		return
	}
	//Создаем и заносим пустое расписание по-умолчанию.
	tempEvenWeek := [][]string{}
	tempOddWeek := [][]string{}
	for i := 0; i < 7; i++ {
		temp := []string{""}
		tempEvenWeek = append(tempEvenWeek, temp)
		tempOddWeek = append(tempOddWeek, temp)
	}
	/////////////////////////////////////////////////////////////////
	///              ДОБАВИТЬ ДАТУ НАЧАЛА ОТСЧЕТА                 ///
	/////////////////////////////////////////////////////////////////
	//Сейчас тут старт у всех в 1 время 5.2
	//Тут callback функция ) удачи/
	err = local_data_base.CreateSchedule(groupName, 5, 2, tempEvenWeek, tempOddWeek)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Ошибка при создании расписания")
		return
	}
	app.InfoLog.Output(2, "Пользователь "+id+" успешно создал группу.")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Группа успешно созданна.")
}

// Обработка запросов на создание группы.
func (app *Application) creatingGroup(w http.ResponseWriter, r *http.Request) {
	//Парсим id из URL
	id := r.URL.Query().Get("id")

	//Парсинг id из string в int64
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", createError("Ошибка при чтении id: "+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, "Ошибка при чтении id.", http.StatusBadRequest)
		return
	}
	////////////////////////////////////////////////////////////////////////////
	////////////////            Зачем это тут?          ////////////////////////
	////////////////////////////////////////////////////////////////////////////
	if !local_data_base.IsAdmin(&intId) {
		trace := fmt.Sprintf("%s\n%s", createError("Пользователь с id "+fmt.Sprint(intId)), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, "Ошибка при чтении id.", http.StatusBadRequest)
		return
	}
	app.InfoLog.Output(2, "Пользователь "+fmt.Sprint(intId)+" успешно подтвердил запрос на создание группы.")
	app.render(w, "create.page.tmpl")
}

// Обработка запросов на удаление группы.
func (app *Application) deleteGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			trace := fmt.Sprintf("%s\n%s", errors.New("error parsing form.\nfunc - deleteGroup; line - 226"), debug.Stack())
			app.ErrorLog.Output(2, trace)
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		groupname := r.Form.Get("groupname")
		err = local_data_base.DeleteGroup(groupname)
		if err != nil {
			trace := fmt.Sprintf("%s\n%s", errors.New("invalid request method.\nfunc - deleteGroup; line - 234"), debug.Stack())
			app.ErrorLog.Output(2, trace)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		app.InfoLog.Output(2, "Группа "+groupname+" - успешно удалена")
		fmt.Fprintf(w, "Группа: %s - успешно удалена", groupname)
	} else {
		trace := fmt.Sprintf("%s\n%s", errors.New("invalid request method, неправильный запрос на удаление группы.\nfunc - deleteGroup; line - 225,243"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
