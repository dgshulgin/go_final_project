package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
)

// Обработчик GET /api/task?id=<идентификатор>
func (server TaskServer) GetTaskById(resp http.ResponseWriter, req *http.Request) {

	id0 := req.URL.Query().Get("id")
	if len(id0) == 0 {
		server.logging(ErrNoId.Error(), nil)
		renderJSON(resp,
			http.StatusBadRequest, dto.Error{Error: ErrNoId.Error()})
		return
	}

	id, err := strconv.Atoi(id0)
	if err != nil {
		msg := errors.Join(ErrInvalidId, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	server.logging(LogGetTaskById, id)

	m, err := server.repo.Get([]uint{uint(id)})
	if err != nil {
		msg := errors.Join(ErrNotExistId, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	// Перенести данные из сущности в DTO
	t := dto.Task{
		Id:      fmt.Sprintf("%d", m[uint(id)].TaskId),
		Date:    m[uint(id)].Date,
		Title:   m[uint(id)].Title,
		Comment: m[uint(id)].Comment,
		Repeat:  m[uint(id)].Repeat,
	}

	// Успех, передать клиенту информацию о задаче
	renderJSON(resp, http.StatusOK, t)
}

// Обработчик GET /api/tasks и редирект на GET /api/tasks?search
func (server TaskServer) GetAllTasks(resp http.ResponseWriter, req *http.Request) {

	queries, ok := req.URL.Query()["search"]
	if ok {
		//поиск по дате ?
		date0, err := time.Parse(formatDateTimeDot, queries[0])
		if err == nil {
			param := date0.Format(formatDateTime)
			http.Redirect(resp, req,
				"/v1/search?date="+param, http.StatusTemporaryRedirect)
			return
		}

		//поиск по тексту
		http.Redirect(resp, req,
			"/v1/search?text="+queries[0], http.StatusTemporaryRedirect)
		return
	}

	//запросить в БД список задач, сортировка по дате ASC
	allTasks, err := server.repo.Get([]uint{})
	if err != nil {
		msg := errors.Join(ErrSelect, err).Error()
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	//Перенести данные в DTO, если задач нет создать пустой DTO
	var tlo dto.TaskList
	if len(allTasks) == 0 {
		tlo.Tasks = []dto.Task{}
	} else {
		for _, task := range allTasks {
			t := dto.Task{
				Id:      fmt.Sprintf("%d", task.TaskId),
				Date:    task.Date,
				Title:   task.Title,
				Comment: task.Comment,
				Repeat:  task.Repeat,
			}
			tlo.Tasks = append(tlo.Tasks, t)
		}
	}

	//Успех, отправить DTO клиенту
	renderJSON(resp, http.StatusOK, tlo)
}
