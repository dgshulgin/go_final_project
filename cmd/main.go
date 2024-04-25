package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dgshulgin/go_final_project/handler"
	"github.com/dgshulgin/go_final_project/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logger.New()
	if err := mainNoExit(log); err != nil {
		log.Fatalf("Аварийное завершение: %s", err.Error())
	}
}

const (
	defaultPort = 7540
)

func mainNoExit(log logrus.FieldLogger) error {
	router, err := handler.Router(log)
	if err != nil {
		return fmt.Errorf("ошибка инициализации маршрутизатора")
	}

	// * Реализуйте возможность определять извне порт при запуске сервера.
	// Порт может быть установлен при запуске сервера через переменную TODO_PORT
	// $ TODO_PORT=8080 go run ./cmd
	// Если переменная TODO_PORT не установлена используется порт по умолчанию 7540.
	Port, ok := checkPortEnv()
	if !ok {
		Port = defaultPort
	}
	log.Printf("Сервер запущен, порт=%d\n", Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", Port), router)
}

func checkPortEnv() (int64, bool) {
	envPort, ok := os.LookupEnv("TODO_PORT")
	if !ok {
		//Переменная TODO_PORT не определена
		return 0, false
	}
	if len(envPort) == 0 {
		//Переменная TODO_PORT определена, но значение не задано
		return 0, false
	}
	eport, err := strconv.ParseInt(envPort, 10, 32)
	if err != nil {
		//Ошибка при конвертации значения
		return 0, false
	}
	return int64(eport), true
}
