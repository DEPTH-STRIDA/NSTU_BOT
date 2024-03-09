package site

import (
	"net/http"
)

// Функция, которая добавляет разные маршруты
func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	//обработка и верификация telegram пользователя
	mux.HandleFunc("/", app.validating)
	//Обработка аякс запросов. Проверяет данные и возвращает успешный код.
	mux.HandleFunc("/validate", app.validate)
	//Главная страница - список групп
	mux.HandleFunc("/home", app.home)
	//Страница с добавлением новой группы.
	mux.HandleFunc("/new-group", app.new_group)
	//Обработка запросов на корректность имени группы.
	mux.HandleFunc("/checkGroupName", app.checkGroupName)
	//Обработка запросов на создание группы.
	mux.HandleFunc("/creatingGroup", app.creatingGroup)
	//Обработка запросов на удаление группы.
	mux.HandleFunc("/deleteGroup", app.deleteGroup)

	//Обработка локальных файлов.
	fileServer := http.FileServer(http.Dir("site/ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
