package task_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	todo "github.com/dgshulgin/go_final_project/internal/pkg/repository/todo"
)

type Handler struct {
	log  logrus.FieldLogger
	repo *todo.Repository
}

func New(log logrus.FieldLogger, repo *todo.Repository) *Handler {
	return &Handler{log, repo}
}

// DTO, задача
type TaskOut struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// DTO, список задач для отправки клиенту
type TaskListOut struct {
	Tasks []TaskOut `json:"tasks"`
}

func (o TaskListOut) Marshal() ([]byte, error) {
	return json.Marshal(o)
}

// DTO, для отправки клиенту сообщения об ошибке
type ErrorOut struct {
	Error string `json:"error"`
}

func (o *ErrorOut) Marshal(e error) ([]byte, error) {
	o.Error = fmt.Sprintf("%s, %s", o.Error, e.Error())
	return json.Marshal(o)
}

// GET
func (h Handler) Get() http.Handler {
	fn := func(resp http.ResponseWriter, req *http.Request) {
		h.log.Printf("обработка запроса GET /api/task")

		//запросить в БД список задач, сортировка по дате ASC
		allTasks, err := h.repo.Get([]uint{})
		if err != nil {
			//напишем в лог
			//отправить ошибку
			return
		}
		//перевалить в DTO, если задач нет DTO пустой
		var tlo TaskListOut
		if len(allTasks) == 0 {
			tlo.Tasks = []TaskOut{}
		} else {
			for _, task := range allTasks {
				t := TaskOut{
					Date:    task.Date,
					Title:   task.Title,
					Comment: task.Comment,
					Repeat:  task.Repeat,
				}
				tlo.Tasks = append(tlo.Tasks, t)
			}
		}

		//отправить DTO клиенту
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		b, err := json.Marshal(tlo)
		if err != nil {
			//закричать
			//отправить ошщибку
			return
		}
		resp.Write(b)
		//FIXME!!!
		h.log.Printf("отправлено, %q", b)
	}
	return http.HandlerFunc(fn)
}
