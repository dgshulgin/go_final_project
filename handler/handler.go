package handler

import (
	"net/http"

	complete "github.com/dgshulgin/go_final_project/handler/complete_task"
	create_task "github.com/dgshulgin/go_final_project/handler/create_task"
	get_tasks "github.com/dgshulgin/go_final_project/handler/get_tasks"
	nextdate_handler "github.com/dgshulgin/go_final_project/handler/nextdate"
	update "github.com/dgshulgin/go_final_project/handler/update_task"
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

	apiNextDate := nextdate_handler.NextDate().ServeHTTP
	router.HandleFunc("/api/nextdate", apiNextDate).Methods(http.MethodGet)

	apiCreateTask := create_task.New(log, repo).Create().ServeHTTP
	router.HandleFunc("/api/task", apiCreateTask).Methods(http.MethodPost)

	apiUpdateTask := update.New(log, repo).Update().ServeHTTP
	router.HandleFunc("/api/task", apiUpdateTask).Methods(http.MethodPut)

	apiGetAllTasks := get_tasks.New(log, repo).GetAll().ServeHTTP
	router.HandleFunc("/api/tasks", apiGetAllTasks).Methods(http.MethodGet)

	apiGetById := get_tasks.New(log, repo).GetById().ServeHTTP
	router.HandleFunc("/api/task", apiGetById).Methods(http.MethodGet)

	apiCompleteTask := complete.New(log, repo).Complete().ServeHTTP
	router.HandleFunc("/api/task/done", apiCompleteTask).Methods(http.MethodPost)

	// я должен быть последним
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(webDir)))

	return router, nil
}
