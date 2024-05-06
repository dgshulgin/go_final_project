package complete

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	task_entity "github.com/dgshulgin/go_final_project/internal/pkg/entity"
	"github.com/dgshulgin/go_final_project/internal/pkg/repository/todo"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	log  logrus.FieldLogger
	repo *todo.Repository
}

// DTO входящего запроса
type TaskIn struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Для отправки клиенту сообщения об ошибке
type ErrorOut struct {
	Error string `json:"error"`
}

func (o *ErrorOut) Marshal(e error) ([]byte, error) {
	o.Error = fmt.Sprintf("%s, %s", o.Error, e.Error())
	return json.Marshal(o)
}

type Ok struct{}

func New(log logrus.FieldLogger, repo *todo.Repository) *Handler {
	return &Handler{log, repo}
}

func (h Handler) Complete() http.Handler {
	fn := func(resp http.ResponseWriter, req *http.Request) {
		h.log.Printf("POST /api/task/done")

		//переваливаем из запроса в DTO
		in := TaskIn{}
		err := json.NewDecoder(req.Body).Decode(&in)
		if err != nil {
			h.log.Errorf("не смог десериализовать json, %s", err.Error())
			resp.WriteHeader(http.StatusBadRequest)
			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			eout := ErrorOut{"не смог десериализовать json"}
			b, _ := eout.Marshal(err)
			resp.Write(b)
			return
		}

		h.log.Printf("параметры запроса {id=%s, date=%s, title=%s, comment=%s, repeat=%s}",
			in.Id, in.Date, in.Title, in.Comment, in.Repeat)

		//валидация DTO
		err = h.validate(&in)
		if err != nil {
			h.log.Errorf("ошибка данных запроса, %s", err.Error())
			resp.WriteHeader(http.StatusBadRequest)
			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			eout := ErrorOut{"невалидный json"}
			b, _ := eout.Marshal(err)
			resp.Write(b)
			return
		}

		h.log.Printf("сохраняем объект как %v", in)

		//переваливаем из DTO в сущность
		//валидация при перевалке
		var t task_entity.Task
		id, _ := strconv.Atoi(in.Id)
		t.TaskId = uint(id)
		t.Date = in.Date
		t.Title = in.Title
		t.Comment = in.Comment
		t.Repeat = in.Repeat

		//сохраняем сущность , для этого у repo должен быть метод Save
		_, err = h.repo.Save(&t)
		if err != nil {
			//не смог сохранить
			//TODO
			return
		}

		//возвращаем id сохраненного объекта
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		ok := Ok{}
		b, _ := json.Marshal(ok)
		resp.Write(b)
		//FIXME!!!
		h.log.Printf("отправлено, %q", b)
	}
	return http.HandlerFunc(fn)
}

func (h *Handler) validate(in *TaskIn) error {
	return nil
}
