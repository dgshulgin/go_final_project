package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
	"github.com/dgshulgin/go_final_project/internal/entity"
)

func (server TaskServer) SearchByText(resp http.ResponseWriter, req *http.Request) {

	//переменная text точно есть и не пустая
	searchText := req.URL.Query().Get("text")

	server.logging(LogSearchText, searchText)

	var t entity.Task
	t.Title = searchText
	t.Comment = searchText

	tasksFound, err := server.repo.Lookup(t)
	if err != nil {
		msg := errors.Join(ErrSearch, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
		return
	}

	// Перенести данные в DTO, если задач нет создать пустой DTO
	var tlo []map[string]string
	if len(tasksFound) == 0 {
		tlo = []map[string]string{}
	} else {

		for _, task := range tasksFound {
			m := make(map[string]string, 5)
			m["id"] = fmt.Sprintf("%d", task.TaskId)
			m["date"] = task.Date
			m["title"] = task.Title
			m["comment"] = task.Comment
			m["repeat"] = task.Repeat
			tlo = append(tlo, m)
		}
	}
	ret := make(map[string][]map[string]string, 1)
	ret["tasks"] = tlo

	// Успех, отправить DTO клиенту
	renderJSON(resp, http.StatusOK, ret)
}

func (server TaskServer) SearchByDate(resp http.ResponseWriter, req *http.Request) {

	//переменная date точно есть и не пустая
	searchDate := req.URL.Query().Get("date")

	server.logging(LogSearchDate, searchDate)

	var t entity.Task
	t.Date = searchDate

	tasksFound, err := server.repo.Lookup(t)
	if err != nil {
		msg := errors.Join(ErrSearch, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
		return
	}

	// //Перенести данные в DTO, если задач нет создать пустой DTO
	var tlo []map[string]string
	if len(tasksFound) == 0 {
		tlo = []map[string]string{}
	} else {

		for _, task := range tasksFound {
			m := make(map[string]string, 5)
			m["id"] = fmt.Sprintf("%d", task.TaskId)
			m["date"] = task.Date
			m["title"] = task.Title
			m["comment"] = task.Comment
			m["repeat"] = task.Repeat
			tlo = append(tlo, m)
		}
	}

	ret := make(map[string][]map[string]string, 1)
	ret["tasks"] = tlo

	// Успех, отправить DTO клиенту
	renderJSON(resp, http.StatusOK, ret)
}
