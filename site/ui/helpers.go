package site

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func (app *Application) render(w http.ResponseWriter, name string) {
	ts, ok := app.TemplateCache[name]
	if !ok {
		fmt.Println(app.TemplateCache)
		app.serverError(w, fmt.Errorf("шаблон %s не существует! ", name))
		return
	}
	err := ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}
