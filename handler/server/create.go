package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
	"github.com/dgshulgin/go_final_project/internal/entity"
	"github.com/dgshulgin/go_final_project/internal/nextdate"
)

// Обработчик POST /api/task/
// Запрос содержит информацию о задаче в формате JSON
// Пример:
//
//	{
//		"date": "20240201",
//		"title": "Подвести итог",
//		"comment": "Мой комментарий",
//		"repeat": "d 5"
//	}
func (server TaskServer) Create(resp http.ResponseWriter, req *http.Request) {

	// Перенести данные из запроса в DTO
	in := dto.Task{}
	err := json.NewDecoder(req.Body).Decode(&in)
	if err != nil {
		msg := fmt.Sprintf("ошибка чтения данных запроса, %s", err.Error())
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	server.log.Printf("создание задачи {id=%s, date=%s, title=%s, comment=%s, repeat=%s}",
		in.Id, in.Date, in.Title, in.Comment, in.Repeat)

	// Валидация DTO
	err = validateOnCreate(&in)
	if err != nil {
		msg := fmt.Sprintf("ошибка валидации данных, %s", err.Error())
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	// Перенести данные из DTO в сущность
	var t entity.Task
	// in.id задачи еще не существует, игнорировать поле t.TaskId
	// t.TaskId
	//
	// in.Date исправлен во время валидации
	t.Date = in.Date
	t.Title = in.Title
	t.Comment = in.Comment
	t.Repeat = in.Repeat

	// Сохранить задачу
	id, err := server.repo.Save(&t)
	if err != nil {
		msg := fmt.Sprintf("ошибка сохранения данных, %s", err.Error())
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
		return
	}

	// Передать клиенту id сохраненной задачи
	renderJSON(resp, http.StatusOK, dto.Id{Id: id})

	in.Id = fmt.Sprintf("%d", id)
	server.log.Printf("задача создана {id=%s, date=%s, title=%s, comment=%s, repeat=%s}",
		in.Id, in.Date, in.Title, in.Comment, in.Repeat)
}

// Проверка валидности данных для входящего запроса, по возможности корректирует
// данные (in.Date)
func validateOnCreate(in *dto.Task) error {

	if len(in.Title) == 0 {
		return fmt.Errorf("поле title не должно быть пустым")
	}

	if len(in.Date) == 0 {
		in.Date = time.Now().Format("20060102")
		return nil
	}

	if strings.EqualFold(in.Date, time.Now().Format("20060102")) {
		return nil
	}

	// if len(in.Repeat) == 0 {
	// 	in.Date = time.Now().Format("20060102")
	// 	return nil
	// }

	//здесь надо проверять что validate возвращает кастомную ошибку и если это так то дополнительно вызывать nextdate для формирования атуальной даты перед сохранением
	err := nextdate.Validate(in.Date, time.Now().Format("20060102"), in.Repeat)
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

	return nil
}
