package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
	"github.com/dgshulgin/go_final_project/internal/entity"
	"github.com/dgshulgin/go_final_project/internal/nextdate"
	"github.com/dgshulgin/go_final_project/services"
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
		msg := errors.Join(ErrBadRequest, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	server.logging(LogCreateTask, in)

	// Валидация DTO
	err = validateOnCreate(&in)
	if err != nil {
		msg := errors.Join(ErrRequestValidation, err).Error()
		server.logging(msg, nil)
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

	// Сохранить сущность
	id, err := server.repo.Save(&t)
	if err != nil {
		msg := errors.Join(ErrSaveData, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
		return
	}

	// Успех, возвратить клиенту id сохраненной задачи
	renderJSON(resp, http.StatusOK, dto.Id{Id: id})
}

// Проверка валидности данных для входящего запроса, коррекция данных, при необходимости.
func validateOnCreate(in *dto.Task) error {

	now := time.Now().Format(services.FormatDateTime)

	if len(in.Title) == 0 {
		return ErrEmptyTitle
	}

	if len(in.Date) == 0 {
		in.Date = now
		return nil
	}

	if strings.EqualFold(in.Date, now) {
		return nil
	}

	// Если validate возвращает ошибку ErrNextDateBeforeNow то дополнительно
	// вызывать nextdate для формирования актуальной даты перед сохранением
	err := nextdate.Validate(in.Date, now, in.Repeat)
	if err != nil {
		if errors.Is(err, nextdate.ErrNextDateBeforeNow) {
			nextDate, err := nextdate.NextDate(in.Date, now, in.Repeat)
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
