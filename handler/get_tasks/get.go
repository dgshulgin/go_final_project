package task_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	Id      string `json:"id"`
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
	//o.Error = fmt.Sprintf("%s", o.Error)
	return json.Marshal(o)
}

func (h Handler) GetById() http.Handler {
	fn := func(resp http.ResponseWriter, req *http.Request) {
		h.log.Printf("обработка запроса GET /api/task")
		req.ParseForm()
		_, ok := req.Form["id"]
		if !ok {
			h.log.Errorf("не указан идентификатор")
			resp.WriteHeader(http.StatusBadRequest)
			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			eout := ErrorOut{"не указан идентификатор"}
			b, _ := json.Marshal(eout)
			h.log.Infof("отправлено %q", b)
			resp.Write(b)
			return
		}

		id := req.URL.Query().Get("id")
		if len(id) == 0 {
			h.log.Errorf("не указан идентификатор")
			resp.WriteHeader(http.StatusBadRequest)
			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			eout := ErrorOut{"не указан идентификатор"}
			b, _ := eout.Marshal(nil)
			resp.Write(b)
			return
		}

		iid, err := strconv.Atoi(id)
		if err != nil {
			h.log.Errorf("некорректный идентификатор %s", id)
			resp.WriteHeader(http.StatusBadRequest)
			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			eout := ErrorOut{"некорректный идентификатор"}
			b, _ := eout.Marshal(err)
			resp.Write(b)
			return
		}
		m, err := h.repo.Get([]uint{uint(iid)})
		if err != nil {
			h.log.Errorf("идентификатор не существует %s", id)
			resp.WriteHeader(http.StatusBadRequest)
			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			eout := ErrorOut{"идентификатор не существует"}
			b, _ := eout.Marshal(err)
			resp.Write(b)
			return
		}
		//h.log.Printf("найдена запись id=%d, data=%v", iid, m[uint(iid)])

		//перевалить в DTO, если задач нет DTO пустой
		//var to TaskOut
		//tlo := make(map[string]TaskOut)
		if len(m) == 0 {
			h.log.Errorf("идентификатор не существует %s", id)
			resp.WriteHeader(http.StatusBadRequest)
			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			eout := ErrorOut{"идентификатор не существует"}
			b, _ := eout.Marshal(err)
			resp.Write(b)
			return
		} else {
			//for _, task := range m {
			t := TaskOut{
				Id:      fmt.Sprintf("%d", m[uint(iid)].TaskId), //task.TaskId),
				Date:    m[uint(iid)].Date,                      //.Date,
				Title:   m[uint(iid)].Title,                     //.Title,
				Comment: m[uint(iid)].Comment,                   //task.Comment,
				Repeat:  m[uint(iid)].Repeat,                    //task.Repeat,
			}
			//tlo[fmt.Sprintf("%d", task.TaskId)] = t
			//отправить DTO клиенту
			resp.WriteHeader(http.StatusOK)
			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			b, err := json.Marshal(t)
			if err != nil {
				h.log.Errorf("не смог сериализовать json, %s", err.Error())
				resp.WriteHeader(http.StatusBadRequest)
				resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
				eout := ErrorOut{"не смог сериализовать json"}
				b, _ := eout.Marshal(err)
				resp.Write(b)
				return
			}
			resp.Write(b)
			//FIXME!!!
			h.log.Printf("отправлено, %q", b)
			//}
		}

	}
	return http.HandlerFunc(fn)
}

// GET
func (h Handler) GetAll() http.Handler {
	fn := func(resp http.ResponseWriter, req *http.Request) {

		h.log.Printf("обработка запроса GET /api/tasks")

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
					Id:      fmt.Sprintf("%d", task.TaskId),
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
