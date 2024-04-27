package handler

import (
	"net/http"

	create_task "github.com/dgshulgin/go_final_project/handler/create_task"
	"github.com/dgshulgin/go_final_project/internal/pkg/repository/todo"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	webDir    = "./web"
	rootRoute = "/"
)

func Router(log logrus.FieldLogger, repo *todo.Repository) (*mux.Router, error) {
	router := mux.NewRouter()

	router.HandleFunc("/api/nextdate", apiNextDate).Methods(http.MethodGet)

	apiCreateTask := create_task.New(log, repo).CreateTask().ServeHTTP
	router.HandleFunc("/api/task", apiCreateTask).Methods(http.MethodPost)

	// я должен быть последним
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(webDir)))

	return router, nil
}
