package task_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	task_entity "github.com/dgshulgin/go_final_project/internal/pkg/entity"
	"github.com/dgshulgin/go_final_project/internal/pkg/nextdate"
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

// DTO входящего запроса
type TaskIn struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Для отправки клиенту идентификатора созданой задачи
type IdOut struct {
	Id uint `json:"id"`
}

func (o IdOut) Marshal() ([]byte, error) {
	return json.Marshal(o)
}

// Для отправки клиенту сообщения об ошибке
type ErrorOut struct {
	Error string `json:"error"`
}

func (o *ErrorOut) Marshal(e error) ([]byte, error) {
	o.Error = fmt.Sprintf("%s, %s", o.Error, e.Error())
	return json.Marshal(o)
}

/*
Валидация DTO, правила
1. поле title обязательное, не должно быть пустым
2. date содержит корректное значение, распознается по формату
3. date = empty, использовать сегодняшнее число
4. date.befor(now) && repeat = empty => now
5. date.befor(now) => nextdate
6. repeat == "mw" => не работаем с форматом
*/
func (h *Handler) validate(in *TaskIn) error {

	if len(in.Title) == 0 {
		return fmt.Errorf("поле title не должно быть пустым")
	}

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

/*
func (h *Handler) validate(in *TaskIn) error {
	//rule 1
	if len(in.Title) == 0 {
		return fmt.Errorf("поле title не может быть пустым")
	}
	// rule 2
	date, err := time.Parse("20060102", in.Date)
	if err != nil {
		return fmt.Errorf("ошибка преобразования даты %s", in.Date)
	}
	//rule 3
	if len(in.Date) == 0 {
		in.Date = time.Now().Format("20060102")
	}
	// rule 5
	if date.Before(time.Now()) {
		if len(in.Repeat) == 0 {
			in.Date = time.Now().Format("20060102")
		}
		nextDate, err := nextdate.NextDate(in.Date, time.Now().Format("20060102"), in.Repeat)
		if err != nil {
			return err
		}
		in.Date = nextDate
	}

	if strings.ContainsAny(in.Repeat, "mw") {
		return fmt.Errorf("поле repeat не поддерживает этот формат")
	}
	return nil
}
*/

// POST
func (h Handler) Create() http.Handler {
	fn := func(resp http.ResponseWriter, req *http.Request) {
		h.log.Printf("обработка запроса POST /api/task")

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

		h.log.Printf("параметры запроса {date=%s, title=%s, comment=%s, repeat=%s}",
			in.Date, in.Title, in.Comment, in.Repeat)

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
		t.Date = in.Date
		t.Title = in.Title
		t.Comment = in.Comment
		t.Repeat = in.Repeat

		//сохраняем сущность , для этого у repo должен быть метод Save
		id, err := h.repo.Save(&t)
		if err != nil {
			//не смог сохранить
			//TODO
			return
		}

		//возвращаем id сохраненного объекта
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
		idout := IdOut{Id: id}
		b, _ := idout.Marshal()
		resp.Write(b)
		//FIXME!!!
		h.log.Printf("отправлено, %q", b)
	}
	return http.HandlerFunc(fn)
}
