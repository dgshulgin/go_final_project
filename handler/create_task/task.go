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
Валидация DTO

	правила валидации

1. title не может быть пустым
2. datе в формате 20060102 и парсится корректно
3. date не задан или содержит пустую сроку = использовать сегодняшнюю дату
4. date = сегодняшняя дата, если repeat не задан или пуст
5. repeat указан в неправильном формате
6. правила w и m возвращают ошибку "не поддерживается", пока
*/
func (h *Handler) validate(in *TaskIn) error {
	//rule 1
	if len(in.Title) == 0 {
		return fmt.Errorf("поле title не может быть пустым")
	}
	//rule 3, rule 4
	if len(in.Date) == 0 || len(in.Repeat) == 0 {
		in.Date = time.Now().Format("20060102")
	}
	// rule 5
	//FIXME вообще говоря это не ошибка и поправимо на уровне хендлера - задача удаляется
	if len(in.Repeat) == 0 {
		return fmt.Errorf("поле repeat не содержит значения")
	}
	// rule 2
	//, rule 6
	// здесь используем NextDate как валидатор
	//значение даты игнорируется так что допустимо использовать текущую дату,
	//какой бы она не была
	_, err := nextdate.NextDate(in.Date,
		time.Now().Format("20060102"), in.Repeat)
	if err != nil {
		return err
	}
	return nil
}

// POST
func (h Handler) CreateTask() http.Handler {
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
