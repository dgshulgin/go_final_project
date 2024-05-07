package complete

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	task_entity "github.com/dgshulgin/go_final_project/internal/pkg/entity"
	"github.com/dgshulgin/go_final_project/internal/pkg/nextdate"
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

		h.log.Printf("POST /api/task/done id=%s", id)

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

		if len(m) == 0 {
			h.log.Errorf("идентификатор не существует %s", id)
			resp.WriteHeader(http.StatusBadRequest)
			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			eout := ErrorOut{"идентификатор не существует"}
			b, _ := eout.Marshal(err)
			resp.Write(b)
			return
		} else {
			repeat := m[uint(iid)].Repeat
			if len(repeat) == 0 {
				err := h.repo.Delete([]uint{uint(iid)})
				if err != nil {
					h.log.Errorf("ошибка при удалении задачи id=%d", iid)
					resp.WriteHeader(http.StatusBadRequest)
					resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
					eout := ErrorOut{"ошибка при удалении задачи"}
					b, _ := eout.Marshal(err)
					resp.Write(b)
					return
				}
				resp.WriteHeader(http.StatusOK)
				resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
				ok := Ok{}
				b, _ := json.Marshal(ok)
				resp.Write(b)
				//FIXME!!!
				h.log.Printf("отправлено, %q", b)
				return
			}
			// возможно повторение, пересчитать nextdate и пересохранить
			var t task_entity.Task
			t.TaskId = m[uint(iid)].TaskId
			nextDate, err := nextdate.NextDate(m[uint(iid)].Date, time.Now().Format("20060102"), m[uint(iid)].Repeat)
			if err != nil {
				//кричать в лог
				//отправить ошибку
				return
			}
			t.Date = nextDate
			t.Title = m[uint(iid)].Title
			t.Comment = m[uint(iid)].Comment
			t.Repeat = m[uint(iid)].Repeat

			//сохраняем сущность , для этого у repo должен быть метод Save
			_, err = h.repo.Save(&t)
			if err != nil {
				//не смог сохранить
				//TODO
				return
			}
			resp.WriteHeader(http.StatusOK)
			resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
			ok := Ok{}
			b, _ := json.Marshal(ok)
			resp.Write(b)
			//FIXME!!!
			h.log.Printf("отправлено, %q", b)
		}
	}
	return http.HandlerFunc(fn)
}
