package site

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gorilla/mux"
)

// Структура приложение.
type Application struct {
	//Логер для ошибок.
	ErrorLog *log.Logger
	//Логер для информационных сообщений.
	InfoLog *log.Logger
	//Карта, хранящая пути к шаблонам.
	TemplateCache map[string]*template.Template
	//Мутекс, обрабатывающий URL пути.
	Router *mux.Router
}

func CreateWebApp(ip, port string) {
	defer print("Тут будет посылание сообщения в канал, о необходимости перезапустить сайт")
	//Log'еры для удобного создания ошибок
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if ip == "" && port == "" {
		trace := fmt.Sprintf("%s\n%s", errors.New("ip и port пустные. Скорее всего переменные среды не найдены"), debug.Stack())
		errorLog.Output(2, trace)
		return
	} else if ip == "" && port != "" {
		trace := fmt.Sprintf("%s\n%s", errors.New("ip пустное, port - не пустое. Скорее всего переменная среды ip не найдена"), debug.Stack())
		errorLog.Output(2, trace)
	} else if ip != "" && port == "" {
		trace := fmt.Sprintf("%s\n%s", errors.New("port пустное, ip - не пустое. Скорее всего переменная среды port не найдена"), debug.Stack())
		errorLog.Output(2, trace)
	}
	addr := flag.String("addr", (ip + ":" + port), "Сетевой адрес веб-сервера")
	flag.Parse()
	templateCache, err := NewTemplateCache("site/ui/html/")
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
		infoLog.Output(2, trace)
		//return
	}
	app := &Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		TemplateCache: templateCache,
	}
	app.Routes()
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
