package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dgshulgin/go_final_project/cmd/config"
	"github.com/dgshulgin/go_final_project/handler/router"
	"github.com/dgshulgin/go_final_project/internal/repository"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	ErrEnvFileNotExist = errors.New("файл окружения отсутствует")
	ErrRepo            = errors.New("ошибка инициализации репозитория")
	ErrRouter          = errors.New("ошибка инициализации маршрутизатора")
	ErrServer          = errors.New("во время работы сервера произошла ошибка")
)

var (
	LogServerFinished = string("сервер завершил работу")
	LogServerStarting = string("запуск сервера %s")
)

func main() {
	log := logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Level:     logrus.DebugLevel,
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf(errors.Join(ErrEnvFileNotExist, err).Error())
	}
	var env config.Config
	host := env.GetEnvAsString("HOST", "localhost")
	port := env.GetEnvAsInt("TODO_PORT", 7540)
	dbname := env.GetEnvAsString("TODO_DBFILE", "scheduler.db")

	repo, err := repository.NewRepository(dbname, &log)
	if err != nil {
		log.Fatalf(errors.Join(ErrRepo, err).Error())
	}
	defer repo.Close()

	router, err := router.NewRouter(&log, repo)
	if err != nil {
		repo.Close()
		log.Errorf(errors.Join(ErrRouter, err).Error())
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}

	log.Printf(LogServerStarting, server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Errorf(errors.Join(ErrServer, err).Error())
	}

	log.Printf(LogServerFinished)

}
