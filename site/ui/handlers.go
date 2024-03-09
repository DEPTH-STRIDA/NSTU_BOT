package site

//Задача данного файла в обработки соответсвующих URL путей.
import (
	"NSTU_NN_BOT/local_data_base"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
)

type Data struct {
	TextSlice []string
}

// Недопустимые слова для названий групп
var unacceptable_words = [...]string{"сегодня", "завтра", "все", "выбор группы", "меню админа", "/open", "/close", "/start"}

// обработка и верификация telegram пользователя
func (app *Application) validating(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "valid.page.tmpl")
}

// Обработка аякс запросов. Проверяет данные и возвращает успешный код.
func (app *Application) validate(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodPost {
		trace := fmt.Sprintf("%s\n%s", errors.New("Invalid request method.\nfunc - validate; line - 26"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, "Invalid request method.\nfunc - validate", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", errors.New("Error parsing form data\nfunc - validate; line - 33"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, "Error parsing form data\nfunc - validate; line - 33", http.StatusBadRequest)
		return
	}

	token := r.Form.Get("token")
	//fmt.Println(token)
	//err = valid.Validate(token, os.Getenv("NSTU_NN_BOT"), 5*time.Minute)
	//if err != nil {
	//fmt.Println(err.Error())
	// Ответ клиенту
	w.WriteHeader(http.StatusOK)
	app.InfoLog.Output(2, "На сайт успешно зашел пользователь. ")
	fmt.Fprintf(w, "Token received successfully: %s", token)
	//} else {
	//// Ответ клиенту
	//	w.WriteHeader(http.StatusOK)
	//	fmt.Fprintf(w, "Token received successfully: %s", token)
	//}
}

// Главная страница - список групп
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	// Извлечение параметра "userID" из строки запроса
	id := r.URL.Query().Get("id")

	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", errors.New("Ошибка конвертации id, возможно пользователь зашел не из telegram.\nfunc - home; line - 62"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		return
	}
	//Если пользователь не является админом, высылаем ему заглушку.
	if local_data_base.IsAdmin(&intId) == false {
		app.render(w, r, "noAdmin.page.tmpl")
		return
	}
	err, groupList := local_data_base.GetGroupsList()
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", errors.New("Ошибка получения списка групп.\nfunc - home; line - 73\n"+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		app.serverError(w, errors.New("Ошибка получения списка групп.\nfunc - home; line - 73\n"+err.Error()))
	}
	textSlice := []string{}
	for _, v := range *groupList {
		if v.Admin == intId {
			textSlice = append(textSlice, v.Name)
		}

	}
	data := Data{textSlice}
	ts, ok := app.TemplateCache["home.page.tmpl"]
	if !ok {
		trace := fmt.Sprintf("%s\n%s", fmt.Errorf("Шаблон %s не существует!", "home.page.tmpl\nfunc - home; line - 87"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		app.serverError(w, fmt.Errorf("Шаблон %s не существует!", "home.page.tmpl"))
		return
	}

	err = ts.Execute(w, &data)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", errors.New("func - home; line - 97\n"+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		app.serverError(w, err)
	}
	app.InfoLog.Output(2, "На главную страницу успешно зашел пользователь. ")
	app.render(w, r, "home.page.tmpl")
}

// Страница с добавлением новой группы.
func (app *Application) new_group(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Output(2, "На страницу создания группы успешно зашел пользователь. ")
	app.render(w, r, "new_group.page.tmpl")
}

// Обработка запросов на корректность имени группы.
func (app *Application) checkGroupName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		trace := fmt.Sprintf("%s\n%s", errors.New("func - checkGroupName; line - 113\nInvalid request method"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", errors.New("func - checkGroupName; line - 119\nИмя использует запрещенные/некорректные символы."), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, "Имя использует запрещенные символы.", http.StatusBadRequest)
		return
	}
	groupName := r.Form.Get("groupName")
	id := r.Form.Get("id")
	err, groupList := local_data_base.GetGroupsList()
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", errors.New("func - checkGroupName; line - 128\nОшибка при получении списка групп."+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		w.WriteHeader(http.StatusBadGateway)
		http.Error(w, "Неизвестная ошибка. Подождите или сообщите в тех. поддержку.", http.StatusBadRequest)
		return
	}
	for _, v := range unacceptable_words {
		if v == groupName {
			trace := fmt.Sprintf("%s\n%s", errors.New("Имя группы: "+groupName+" - недопустимо"), debug.Stack())
			app.ErrorLog.Output(2, trace)
			w.WriteHeader(http.StatusBadGateway)
			http.Error(w, "Неизвестная ошибка. Подождите или сообщите в тех. поддержку.", http.StatusBadRequest)
			return
		}
	}
	for _, v := range *groupList {
		if v.Name == groupName {
			trace := fmt.Sprintf("%s\n%s", errors.New("Имя группы: "+groupName+" - занято"), debug.Stack())
			app.ErrorLog.Output(2, trace)
			w.WriteHeader(http.StatusBadGateway)
			http.Error(w, "Неизвестная ошибка. Подождите или сообщите в тех. поддержку.", http.StatusBadRequest)
			return
		}
	}
	// Конвертация строки в int64
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", errors.New("Ошибка при конвертации id; func - checkGroupName; line - 158\n"+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Ошибка при конвертации id")
		return
	}
	err = local_data_base.CreateGroup(groupName, intId)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", errors.New("Ошибка при добавлении группы; func - checkGroupName; line - 166\n"+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Ошибка при добавлении группы")
		return
	}
	////////////////////////////////////////////////////////////////////////////////
	//////////////////////////       Создание расписания            ///////////////
	//////////////////////////////////////////////////////////////////////////////
	//Сейчас заглушка
	tempEvenWeek := [][]string{}
	tempOddWeek := [][]string{}
	for i := 0; i < 7; i++ {
		temp := []string{"Не заполнено"}
		tempEvenWeek = append(tempEvenWeek, temp)
		tempOddWeek = append(tempOddWeek, temp)
	}

	////////////////////////////////////////////////////////////////////////////////
	//Тут callback функция ) удачи
	err = local_data_base.CreateSchedule(groupName, func() int {
		todayDay := time.Now().Day()
		if todayDay == 1 {
			return 1
		}
		return todayDay - 1
	}(), int(time.Now().Month()), tempEvenWeek, tempOddWeek)
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
	id := r.URL.Query().Get("id")

	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", errors.New("Ошибка при чтении id; func - creatingGroup; line - 209\n"+err.Error()), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, "Ошибка при чтении id.", http.StatusBadRequest)
		return
	}
	////////////////////////////////////////////////////////////////////////////
	////////////////            Зачем это тут?          ////////////////////////
	////////////////////////////////////////////////////////////////////////////
	app.InfoLog.Output(2, "Пользователь "+fmt.Sprint(intId)+" успешно подтвердил запрос на создание группы.")
	app.render(w, r, "create.page.tmpl")
}

// Обработка запросов на удаление группы.
func (app *Application) deleteGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			trace := fmt.Sprintf("%s\n%s", errors.New("Error parsing form.\nfunc - deleteGroup; line - 226"), debug.Stack())
			app.ErrorLog.Output(2, trace)
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		groupname := r.Form.Get("groupname")
		err = local_data_base.DeleteGroup(groupname)
		if err != nil {
			trace := fmt.Sprintf("%s\n%s", errors.New("Invalid request method.\nfunc - deleteGroup; line - 234"), debug.Stack())
			app.ErrorLog.Output(2, trace)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		app.InfoLog.Output(2, "Группа "+groupname+" - успешно удалена")
		fmt.Fprintf(w, "Группа: %s - успешно удалена", groupname)
	} else {
		trace := fmt.Sprintf("%s\n%s", errors.New("Invalid request method, неправильный запрос на удаление группы.\nfunc - deleteGroup; line - 225,243"), debug.Stack())
		app.ErrorLog.Output(2, trace)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
