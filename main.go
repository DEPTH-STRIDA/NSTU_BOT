// Пакет является точкой входа в программу.
// Тут запускается бот и веб приложение.
package main

import (
	site "NSTU_NN_BOT/site/ui"

	bot "NSTU_NN_BOT/telegram_bot"
	"os"
)

func main() {
	//Запуск Телеграмм бота.
	go bot.CreateTgBot(os.Getenv("NSTU_NN_BOT"))
	//Запуск веб приложения.
	site.CreateWebApp(os.Getenv("APP_IP"), os.Getenv("APP_PORT"))
}
