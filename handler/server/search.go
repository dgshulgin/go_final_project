package server

import (
	"fmt"
	"net/http"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
	"github.com/dgshulgin/go_final_project/internal/entity"
)

func (server TaskServer) SearchByText(resp http.ResponseWriter, req *http.Request) {

	//переменная text точно есть и не пустая
	searchText := req.URL.Query().Get("text")

	server.log.Printf("поиск по тексту заголовка или комментария \"%s\"", searchText)

	var t entity.Task
	t.Title = searchText
	t.Comment = searchText

	tasksFound, err := server.repo.Lookup(t)
	if err != nil {
		msg := fmt.Sprintf("ошибка во время поиска данных, %s", err.Error())
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
		return
	}

	// //Перенести данные в DTO, если задач нет создать пустой DTO
	// var tlo dto.TaskList
	// if len(tasksFound) == 0 {
	// 	tlo.Tasks = []dto.Task{}
	// } else {
	// 	for _, task := range tasksFound {
	// 		t := dto.Task{
	// 			Id:      fmt.Sprintf("%d", task.TaskId),
	// 			Date:    task.Date,
	// 			Title:   task.Title,
	// 			Comment: task.Comment,
	// 			Repeat:  task.Repeat,
	// 		}
	// 		tlo.Tasks = append(tlo.Tasks, t)
	// 	}
	// }
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
	//var ret map[string][]map[string]string
	ret := make(map[string][]map[string]string, 1)
	ret["tasks"] = tlo

	//отправить DTO клиенту
	renderJSON(resp, http.StatusOK, ret)

	server.log.Printf("отправлены задачи %q", ret)
}

func (server TaskServer) SearchByDate(resp http.ResponseWriter, req *http.Request) {

	//переменная date точно есть и не пустая
	searchDate := req.URL.Query().Get("date")

	server.log.Printf("поиск по дате %s", searchDate)

	var t entity.Task
	t.Date = searchDate

	tasksFound, err := server.repo.Lookup(t)
	if err != nil {
		msg := fmt.Sprintf("ошибка во время поиска данных, %s", err.Error())
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
		return
	}

	// //Перенести данные в DTO, если задач нет создать пустой DTO
	// var tlo dto.TaskList
	// if len(tasksFound) == 0 {
	// 	tlo.Tasks = []dto.Task{}
	// } else {
	// 	for _, task := range tasksFound {
	// 		t := dto.Task{
	// 			Id:      fmt.Sprintf("%d", task.TaskId),
	// 			Date:    task.Date,
	// 			Title:   task.Title,
	// 			Comment: task.Comment,
	// 			Repeat:  task.Repeat,
	// 		}
	// 		tlo.Tasks = append(tlo.Tasks, t)
	// 	}
	// }

	// //отправить DTO клиенту
	// renderJSON(resp, http.StatusOK, tlo)

	// server.log.Printf("отправлены задачи %q", tlo)

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
	//var ret map[string][]map[string]string
	ret := make(map[string][]map[string]string, 1)
	ret["tasks"] = tlo

	//отправить DTO клиенту
	renderJSON(resp, http.StatusOK, ret)

	server.log.Printf("отправлены задачи %q", ret)
}
