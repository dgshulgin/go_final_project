package delete

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func (h Handler) Delete() http.Handler {
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

		h.log.Printf("DELETE /api/task/done id=%s", id)

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

		err = h.repo.Delete([]uint{uint(iid)})
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
		b, _ := json.Marshal(Ok{})
		resp.Write(b)
		//FIXME!!!
		h.log.Printf("отправлено, %q", b)
	}
	return http.HandlerFunc(fn)
}
