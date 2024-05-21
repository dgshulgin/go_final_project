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

	server := server.NewTaskServer(log, repo)

	r := mux.NewRouter().StrictSlash(true)
	r.Use(MiddlewareLogger(log))

	r.HandleFunc("/api/nextdate", server.NextDate).Methods(http.MethodGet)
	r.HandleFunc("/api/signin", server.Authenticate).Methods(http.MethodPost)

	taskListRouter := r.PathPrefix("/api/tasks").Subrouter()
	taskListRouter.Methods(http.MethodGet).HandlerFunc(server.GetAllTasks)
	taskListRouter.Use(server.MiddlewareCheckUserAuth())

	taskRouter := r.PathPrefix("/api/task").Subrouter()
	taskRouter.HandleFunc("/done", server.Complete).Methods(http.MethodPost)
	taskRouter.Methods(http.MethodPost).HandlerFunc(server.Create)
	taskRouter.Methods(http.MethodPut).HandlerFunc(server.Update)
	taskRouter.Methods(http.MethodDelete).HandlerFunc(server.Delete)
	taskRouter.Methods(http.MethodGet).HandlerFunc(server.GetTaskById)
	taskRouter.Use(server.MiddlewareCheckUserAuth())

	//taskRouter.HandleFunc("/done", server.Complete).Methods(http.MethodPost)
	//taskRouter.PathPrefix("/done").Methods(http.MethodPost).HandlerFunc(server.Complete)

	searchRouter := r.PathPrefix("/v1/search").Subrouter()
	searchRouter.Methods(http.MethodGet).Queries("text", "").HandlerFunc(server.SearchByText)
	searchRouter.Methods(http.MethodGet).Queries("date", "").HandlerFunc(server.SearchByDate)

	rootRouter := r.PathPrefix("/").Subrouter()
	rootRouter.Methods(http.MethodGet).HandlerFunc(http.FileServer(http.Dir(webDir)).ServeHTTP)

	//router.HandleFunc("/api/task", server.Create).Methods(http.MethodPost)
	//router.HandleFunc("/api/task", server.Update).Methods(http.MethodPut)
	//router.HandleFunc("/api/task", server.Delete).Methods(http.MethodDelete)
	//router.HandleFunc("/api/task/done", server.Complete).Methods(http.MethodPost)
	//router.HandleFunc("/api/task", server.GetTaskById).Methods(http.MethodGet)

	//router.HandleFunc("/api/tasks", server.GetAllTasks).Methods(http.MethodGet)

	// router.HandleFunc("/api/search", server.SearchByText).Methods(http.MethodGet).Queries("text", "")
	// router.HandleFunc("/api/search", server.SearchByDate).Methods(http.MethodGet).Queries("date", "")

	//return router, nil
	return r, nil
}

func MiddlewareLogger(log logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			log.Printf("обработка %s - %s (%s)", req.Method, req.URL.Path, req.RemoteAddr)
			next.ServeHTTP(resp, req)
		})
	}
}
