package server

import (
	"encoding/json"
	"errors"
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
		msg := errors.Join(ErrBadRequest, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	server.logging(LogUpdateTask, in)

	//валидация DTO
	err = validateOnUpdate(&in, server.repo)
	if err != nil {
		msg := errors.Join(ErrRequestValidation, err).Error()
		server.logging(msg, nil)
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
		msg := errors.Join(ErrSaveData, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
		return
	}

	// Успех, возвратить пустой JSON
	renderJSON(resp, http.StatusOK, dto.Ok{})
}

// Проверка валидности данных для входящего запроса, коррекция данных, при необходимости.
func validateOnUpdate(in *dto.Task, repo *repository.Repository) error {

	if len(in.Title) == 0 {
		return ErrEmptyTitle
	}

	if len(in.Id) == 0 {
		return ErrNoId
	}

	id0, err := strconv.Atoi(in.Id)
	if err != nil {
		return errors.Join(ErrInvalidId, err)
	}

	_, err = repo.Get([]uint{uint(id0)})
	if err != nil {
		return ErrNotExistId
	}

	now := time.Now().Format(formatDateTime)

	// Если validate возвращает ошибку ErrNextDateBeforeNow то дополнительно
	// вызывать nextdate для формирования актуальной даты перед сохранением
	err = nextdate.Validate(in.Date, now, in.Repeat)
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
