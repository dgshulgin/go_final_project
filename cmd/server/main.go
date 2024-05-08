package main

import (
	"fmt"
	"net/http"
	"os"

	router "github.com/dgshulgin/go_final_project/handler/router"
	"github.com/dgshulgin/go_final_project/internal/logger"
	"github.com/dgshulgin/go_final_project/internal/repository"
	"github.com/sirupsen/logrus"
)

const (
	defaultPort   = "7540"
	defaultDbName = "scheduler.db"
)

func main() {
	log := logger.New()
	log.Level = logrus.DebugLevel
	if err := mainNoExit(log); err != nil {
		log.Fatalf("Аварийное завершение, %s", err.Error())
	}
}

func mainNoExit(log logrus.FieldLogger) error {
	// * Реализуйте возможность определять путь к файлу базы данных через
	//переменную окружения.
	// $ TODO_DBFILE=./database.db go run ./cmd
	//Если переменная TODO_DBFILE не установлена используется имя файла
	//по умолчанию scheduler.db. Файл БД по умолчанию должен находиться
	//рядом с исполняемым файлом, в противном случае файл БД будет создан
	//заново.
	//dbName, ok := checkDbEnv()
	dbName, ok := readEnv("TODO_DBFILE")
	if !ok {
		log.Infof("переменная окружения TODO_DBFILE не определена, БД по умолчанию %s", defaultDbName)
		dbName = defaultDbName
	}

	// * Реализуйте возможность определять извне порт при запуске сервера.
	// $ TODO_PORT=8080 go run ./cmd
	// Если переменная TODO_PORT не установлена используется порт по умолчанию 7540.
	//Port, ok := checkPortEnv()
	port, ok := readEnv("TODO_PORT")
	if !ok {
		log.Infof("переменная окружения TODO_PORT не определена, порт по умолчанию %s", defaultPort)
		port = defaultPort
	}

	repo, err := repository.NewRepository(dbName, log)
	if err != nil {
		return fmt.Errorf("ошибка инициализации репозитория, %w", err)
	}
	defer repo.Close()

	router, err := router.NewRouter(log, repo)
	if err != nil {
		return fmt.Errorf("ошибка инициализации маршрутизатора, %w", err)
	}

	srv := &http.Server{
		Handler: router,
		Addr:    "localhost:" + port,
	}

	log.Infof("Запуск сервера, адрес %s", srv.Addr)
	return srv.ListenAndServe()
}

// Возвращает значение переменной окружения
func readEnv(name string) (string, bool) {
	env, ok := os.LookupEnv(name)
	if !ok {
		//Переменная окружения не определена
		return "", false
	}
	if len(env) == 0 {
		//Переменная окружения определена, но значение не задано
		return "", false
	}
	return env, true
}
