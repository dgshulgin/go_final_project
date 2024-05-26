package server

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
	task_entity "github.com/dgshulgin/go_final_project/internal/entity"
	"github.com/dgshulgin/go_final_project/internal/nextdate"
	"github.com/dgshulgin/go_final_project/services"
)

// Обработчик POST /api/task/done?id=<идентификатор>
func (server TaskServer) Complete(resp http.ResponseWriter, req *http.Request) {

	id0 := req.URL.Query().Get("id")
	if len(id0) == 0 {
		msg := ErrNoId.Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	id, err := strconv.Atoi(id0)
	if err != nil {
		msg := errors.Join(ErrInvalidId, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	m, err := server.repo.Get([]uint{uint(id)})
	if err != nil {
		msg := errors.Join(ErrNotExistId, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	// Если повторений нет, задача считается завершенной
	repeat := m[uint(id)].Repeat
	if len(repeat) == 0 {

		server.logging(LogCompleteTask, id)

		err := server.repo.Delete([]uint{uint(id)})
		if err != nil {
			msg := errors.Join(ErrDelete, err).Error()
			server.logging(msg, nil)
			renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
			return
		}

		// Успех, задача завершена и удалена
		renderJSON(resp, http.StatusOK, dto.Ok{})
		return
	}

	server.logging(LogUpdateTask, m[uint(id)].TaskId)

	// возможно повторение, пересчитать nextdate и пересохранить
	now := time.Now().Format(services.FormatDateTime)
	err = nextdate.Validate(m[uint(id)].Date, now, m[uint(id)].Repeat)
	if err != nil {
		server.logging(err.Error(), nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: err.Error()})
		return
	}

	nextDate, err := nextdate.NextDate(m[uint(id)].Date, now, m[uint(id)].Repeat)
	if err != nil {
		server.logging(err.Error(), nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: err.Error()})
		return
	}

	// Перевалим в DTO
	var t task_entity.Task
	t.TaskId = m[uint(id)].TaskId
	t.Date = nextDate
	t.Title = m[uint(id)].Title
	t.Comment = m[uint(id)].Comment
	t.Repeat = m[uint(id)].Repeat

	//Сохранить изменения в задаче
	_, err = server.repo.Save(&t)
	if err != nil {
		msg := errors.Join(ErrSaveData, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
		return
	}

	// Успех, возвращаем пустой JSON
	renderJSON(resp, http.StatusOK, dto.Ok{})
}
