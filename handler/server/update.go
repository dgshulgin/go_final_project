package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
	"github.com/dgshulgin/go_final_project/internal/entity"
	"github.com/dgshulgin/go_final_project/internal/nextdate"
	"github.com/dgshulgin/go_final_project/internal/repository"
)

// Обработчик PUT /api/task
// Запрос содержит информацию о задаче в формате JSON
//
//	{
//	    "id": "185",
//	    "date": "20240201",
//	    "title": "Подвести итог",
//	    "comment": "",
//	    "repeat": ""
//	}
func (server TaskServer) Update(resp http.ResponseWriter, req *http.Request) {

	// Перенести данные из запроса в DTO
	in := dto.Task{}
	err := json.NewDecoder(req.Body).Decode(&in)
	if err != nil {
		msg := fmt.Sprintf("Update: ошибка чтения данных запроса, %s", err.Error())
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	server.log.Printf(
		"обновление задачи {id=%s, date=%s, title=%s, comment=%s, repeat=%s}",
		in.Id, in.Date, in.Title, in.Comment, in.Repeat)

	//валидация DTO
	err = validateOnUpdate(&in, server.repo)
	if err != nil {
		msg := fmt.Sprintf("ошибка валидации данных, %s", err.Error())
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	// Перенести данные из DTO в сущность
	var t entity.Task
	id0, _ := strconv.Atoi(in.Id)
	t.TaskId = uint(id0)
	t.Date = in.Date
	t.Title = in.Title
	t.Comment = in.Comment
	t.Repeat = in.Repeat

	// Сохранить задачу
	_, err = server.repo.Save(&t)
	if err != nil {
		msg := fmt.Sprintf("ошибка сохранения данных, %s", err.Error())
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
		return
	}

	// Информировать клиента об успешном обновлении
	renderJSON(resp, http.StatusOK, dto.Ok{})

	server.log.Printf("задача обновлена {id=%s, date=%s, title=%s, comment=%s, repeat=%s}",
		in.Id, in.Date, in.Title, in.Comment, in.Repeat)
}

// Проверка валидности данных для входящего запроса, по возможности корректирует
// данные (in.Date)
func validateOnUpdate(in *dto.Task, repo *repository.Repository) error {

	if len(in.Title) == 0 {
		return fmt.Errorf("поле title не должно быть пустым")
	}

	if len(in.Id) == 0 {
		return fmt.Errorf("некорректное значение поля id=%s", in.Id)
	}

	iid, err := strconv.Atoi(in.Id)
	if err != nil {
		return fmt.Errorf("некорректное значение поля id=%s", in.Id)
	}

	_, err = repo.Get([]uint{uint(iid)})
	if err != nil {
		return fmt.Errorf("id=%d не существует", iid)
	}

	//здесь надо проверять что validate возвращает кастомную ошибку и если это так то дополнительно вызывать nextdate для формирования атуальной даты перед сохранением
	err = nextdate.Validate(in.Date, time.Now().Format("20060102"), in.Repeat)
	if err != nil {
		if errors.Is(err, nextdate.ErrNextDateBeforeNow) {
			nextDate, err := nextdate.NextDate(in.Date, time.Now().Format("20060102"), in.Repeat)
			if err != nil {
				return err
			}
			in.Date = nextDate
			return nil
		}
		return err
	}

	// nextDate, err := nextdate.NextDate(in.Date, time.Now().Format("20060102"), in.Repeat)
	// if err != nil {
	// 	return err
	// }
	// in.Date = nextDate

	return nil
}
