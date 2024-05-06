package update

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

// Для отправки клиенту идентификатора созданой задачи
type IdOut struct {
	Id uint `json:"id"`
}

func New(log logrus.FieldLogger, repo *todo.Repository) *Handler {
	return &Handler{log, repo}
}

func (h Handler) Update() http.Handler {
	fn := func(resp http.ResponseWriter, req *http.Request) {
		h.log.Printf("PUT /api/task")

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

	if len(in.Title) == 0 {
		return fmt.Errorf("поле title не должно быть пустым")
	}

	if len(in.Id) == 0 {
		return fmt.Errorf("некорректное значение поля id")
	}

	iid, err := strconv.Atoi(in.Id)
	if err != nil {
		return fmt.Errorf("некорректное значение поля id")
	}

	m, err := h.repo.Get([]uint{uint(iid)})
	if err != nil {
		return fmt.Errorf("указанный id не существует")
	}
	h.log.Printf("найдена запись id=%d, data=%v", iid, m[uint(iid)])

	// if strings.ContainsAny(in.Repeat, "wm") {
	// 	return fmt.Errorf("формат w или m не поддерживается")
	// }
	// if len(in.Date) == 0 {
	// 	in.Date = time.Now().Format("20060102")
	// 	return nil
	// }

	// startDate, err := time.Parse("20060102", in.Date)
	// if err != nil {
	// 	return fmt.Errorf("некорректное значение поля date")
	// }

	// if startDate.Before(time.Now()) {
	// 	//in.Date = time.Now().AddDate(0, 0, 1).Format("20060102")
	// 	if len(in.Repeat) > 0 {
	// 		in.Date = time.Now().Format("20060102")
	// 		return nil
	// 	}
	// 	// else {
	// 	// 	nextDate, err := nextdate.NextDate(in.Date, time.Now().Format("20060102"), in.Repeat)
	// 	// 	if err != nil {
	// 	// 		return err
	// 	// 	}
	// 	// 	in.Date = nextDate
	// 	// }
	// }

	// //проверить валиадность самого правила посторения
	// if len(in.Repeat) > 0 {
	nextDate, err := nextdate.NextDate(in.Date, time.Now().Format("20060102"), in.Repeat)
	if err != nil {
		return err
	}
	in.Date = nextDate
	// }

	return nil
}
