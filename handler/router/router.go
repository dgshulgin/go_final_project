package router

import (
	"net/http"

	"github.com/dgshulgin/go_final_project/handler/server"
	"github.com/dgshulgin/go_final_project/internal/repository"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	webDir = "./web"
)

func NewRouter(log logrus.FieldLogger, repo *repository.Repository) (*mux.Router, error) {
	router := mux.NewRouter()
	router.StrictSlash(true)

	server := server.NewTaskServer(log, repo)
	router.HandleFunc("/api/nextdate", server.NextDate).Methods(http.MethodGet)
	router.HandleFunc("/api/task", server.Create).Methods(http.MethodPost)
	router.HandleFunc("/api/task", server.Update).Methods(http.MethodPut)
	router.HandleFunc("/api/tasks", server.GetAll).Methods(http.MethodGet)
	router.HandleFunc("/api/task", server.GetById).Methods(http.MethodGet)
	router.HandleFunc("/api/task/done", server.Complete).Methods(http.MethodPost)
	router.HandleFunc("/api/task", server.Delete).Methods(http.MethodDelete)

	// я должен быть последним
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(webDir)))

	return router, nil
}
