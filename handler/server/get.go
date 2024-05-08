package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
)

// Обработчик GET /api/task?id=<идентификатор>
func (server TaskServer) GetById(resp http.ResponseWriter, req *http.Request) {

	id0 := req.URL.Query().Get("id")
	if len(id0) == 0 {
		msg := "GetById: не указан идентификатор задачи"
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	server.log.Printf("запрос информации о задаче id=%s", id0)

	id, err := strconv.Atoi(id0)
	if err != nil {
		msg := fmt.Sprintf("некорректный идентификатор задачи %d", id)
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	m, err := server.repo.Get([]uint{uint(id)})
	if err != nil {
		msg := fmt.Sprintf("идентификатор не существует, %d", id)
		server.log.Errorf(msg)
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

	// Передать клиенту информацию о задаче
	renderJSON(resp, http.StatusOK, t)

	server.log.Infof("отправлена задача %q", t)
}

// Обработчик GET /api/tasks
func (server TaskServer) GetAll(resp http.ResponseWriter, req *http.Request) {

	server.log.Printf("запрос информации о всех задачах")

	//запросить в БД список задач, сортировка по дате ASC
	allTasks, err := server.repo.Get([]uint{})
	if err != nil {
		msg := fmt.Sprintf("ошибка при запросе к БД, %s", err.Error())
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

	//отправить DTO клиенту
	renderJSON(resp, http.StatusOK, tlo)

	server.log.Infof("отправлены задачи %q", tlo)
}
