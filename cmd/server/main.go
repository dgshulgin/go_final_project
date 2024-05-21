package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgshulgin/go_final_project/cmd/config"
	"github.com/dgshulgin/go_final_project/handler/router"
	"github.com/dgshulgin/go_final_project/internal/repository"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Level:     logrus.DebugLevel,
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("файл окружения отсутствует")
	}
	var env config.Config
	host := env.GetEnvAsString("HOST", "localhost")
	port := env.GetEnvAsInt("TODO_PORT", 7540)
	dbname := env.GetEnvAsString("TODO_DBFILE", "scheduler.db")

	repo, err := repository.NewRepository(dbname, &log)
	if err != nil {
		log.Fatalf("ошибка инициализации репозитория, %s", err.Error())
	}
	defer repo.Close()

	router, err := router.NewRouter(&log, repo)
	if err != nil {
		repo.Close()
		log.Errorf("ошибка инициализации маршрутизатора, %s", err.Error())
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}

	log.Printf("запуск сервера %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Errorf("во время работы сервера произошла ошибка %s", err.Error())
	}

	log.Printf("сервер завершил работу")

}

// const (
// 	defaultPort   = "7540"
// 	defaultDbName = "scheduler.db"
// )

// func main() {
// 	log := logger.New()
// 	log.Level = logrus.DebugLevel
// 	if err := mainNoExit(log); err != nil {
// 		log.Fatalf("Аварийное завершение, %s", err.Error())
// 	}
// }

// func mainNoExit(log logrus.FieldLogger) error {
// 	// * Реализуйте возможность определять путь к файлу базы данных через
// 	//переменную окружения.
// 	// $ TODO_DBFILE=./database.db go run ./cmd
// 	//Если переменная TODO_DBFILE не установлена используется имя файла
// 	//по умолчанию scheduler.db. Файл БД по умолчанию должен находиться
// 	//рядом с исполняемым файлом, в противном случае файл БД будет создан
// 	//заново.
// 	//dbName, ok := checkDbEnv()
// 	dbName, ok := readEnv("TODO_DBFILE")
// 	if !ok {
// 		log.Infof("переменная окружения TODO_DBFILE не определена, БД по умолчанию %s", defaultDbName)
// 		dbName = defaultDbName
// 	}

// 	// * Реализуйте возможность определять извне порт при запуске сервера.
// 	// $ TODO_PORT=8080 go run ./cmd
// 	// Если переменная TODO_PORT не установлена используется порт по умолчанию 7540.
// 	//Port, ok := checkPortEnv()
// 	port, ok := readEnv("TODO_PORT")
// 	if !ok {
// 		log.Infof("переменная окружения TODO_PORT не определена, порт по умолчанию %s", defaultPort)
// 		port = defaultPort
// 	}

// 	repo, err := repository.NewRepository(dbName, log)
// 	if err != nil {
// 		return fmt.Errorf("ошибка инициализации репозитория, %w", err)
// 	}
// 	defer repo.Close()

// 	router, err := router.NewRouter(log, repo)
// 	if err != nil {
// 		return fmt.Errorf("ошибка инициализации маршрутизатора, %w", err)
// 	}

// 	srv := &http.Server{
// 		Handler: router,
// 		Addr:    "localhost:" + port,
// 	}

// 	log.Infof("Запуск сервера, адрес %s", srv.Addr)
// 	return srv.ListenAndServe()
// }

// // Возвращает значение переменной окружения
// func readEnv(name string) (string, bool) {
// 	env, ok := os.LookupEnv(name)
// 	if !ok {
// 		//Переменная окружения не определена
// 		return "", false
// 	}
// 	if len(env) == 0 {
// 		//Переменная окружения определена, но значение не задано
// 		return "", false
// 	}
// 	return env, true
// }
