package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dgshulgin/go_final_project/handler"
	todo "github.com/dgshulgin/go_final_project/internal/pkg/repository/todo"
	"github.com/dgshulgin/go_final_project/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logger.New()
	if err := mainNoExit(log); err != nil {
		log.Fatalf("Аварийное завершение, %s", err.Error())
	}
}

const (
	defaultPort   = 7540
	defaultDbName = "scheduler.db"
)

func mainNoExit(log logrus.FieldLogger) error {
	// * Реализуйте возможность определять путь к файлу базы данных через
	//переменную окружения.
	// $ TODO_DBFILE=./database.db go run ./cmd
	//Если переменная TODO_DBFILE не установлена используется имя файла
	//по умолчанию scheduler.db. Файл БД по умолчанию должен находиться
	//рядом с исполняемым файлом, в противном случае файл БД будет создан
	//заново.
	dbName, ok := checkDbEnv()
	if !ok {
		dbName = defaultDbName
	}
	repo, err := todo.NewRepository(dbName, log)
	if err != nil {
		return fmt.Errorf("ошибка репозитория, %w", err)
	}
	defer repo.Close()

	router, err := handler.Router(log, repo)
	if err != nil {
		return fmt.Errorf("ошибка инициализации маршрутизатора, %w", err)
	}

	// * Реализуйте возможность определять извне порт при запуске сервера.
	// $ TODO_PORT=8080 go run ./cmd
	// Если переменная TODO_PORT не установлена используется порт по умолчанию 7540.
	Port, ok := checkPortEnv()
	if !ok {
		Port = defaultPort
	}

	srv := &http.Server{
		Handler: router,
		Addr:    "localhost:" + fmt.Sprintf("%d", Port),
	}

	log.Printf("Сервер запущен, порт=%d\n", Port)
	return srv.ListenAndServe()
}

func checkDbEnv() (string, bool) {
	envDBF, ok := os.LookupEnv("TODO_DBFILE")
	if !ok { //Переменная TODO_DBFILE не определена
		return "", false
	}
	if len(envDBF) == 0 { //Переменная TODO_DBF определена, но значение не задано
		return "", false
	}
	return envDBF, true
}

func checkPortEnv() (int64, bool) {
	envPort, ok := os.LookupEnv("TODO_PORT")
	if !ok { //Переменная TODO_PORT не определена
		return 0, false
	}
	if len(envPort) == 0 { //Переменная TODO_PORT определена, но значение не задано
		return 0, false
	}
	eport, err := strconv.ParseInt(envPort, 10, 32)
	if err != nil { //Ошибка при конвертации значения
		return 0, false
	}
	return int64(eport), true
}
