package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
	task_entity "github.com/dgshulgin/go_final_project/internal/entity"
	"github.com/dgshulgin/go_final_project/internal/nextdate"
)

// Обработчик POST /api/task/done?id=<идентификатор>
func (server TaskServer) Complete(resp http.ResponseWriter, req *http.Request) {

	id0 := req.URL.Query().Get("id")
	if len(id0) == 0 {
		msg := "Complete: не указан идентификатор"
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	server.log.Printf("запрос на завершение задачи id=%s", id0)

	id, err := strconv.Atoi(id0)
	if err != nil {
		msg := fmt.Sprintf("некорректный идентификатор %s", id0)
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	m, err := server.repo.Get([]uint{uint(id)})
	if err != nil {
		msg := fmt.Sprintf("идентификатор не существует %d", id)
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	// Если повторений нет, задача считается завершенной
	repeat := m[uint(id)].Repeat
	if len(repeat) == 0 {
		err := server.repo.Delete([]uint{uint(id)})
		if err != nil {
			msg := fmt.Sprintf("ошибка при удалении задачи id=%d", id)
			server.log.Errorf(msg)
			renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
			return
		}
		renderJSON(resp, http.StatusOK, dto.Ok{})
		server.log.Infof("условия повторения отсутствуют, задача завершена, удалена %q", m[uint(id)])
		return
	}

	// возможно повторение, пересчитать nextdate и пересохранить
	var t task_entity.Task
	t.TaskId = m[uint(id)].TaskId
	nextDate, err := nextdate.NextDate(m[uint(id)].Date, time.Now().Format("20060102"), m[uint(id)].Repeat)
	if err != nil {
		server.log.Errorf(err.Error())
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: err.Error()})
		return
	}
	t.Date = nextDate
	t.Title = m[uint(id)].Title
	t.Comment = m[uint(id)].Comment
	t.Repeat = m[uint(id)].Repeat

	//Сохранить изменения в задаче
	_, err = server.repo.Save(&t)
	if err != nil {
		msg := fmt.Sprintf("ошибка сохранения данных, %s", err.Error())
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusInternalServerError, dto.Error{Error: msg})
		return
	}
	renderJSON(resp, http.StatusOK, dto.Ok{})

	server.log.Infof("задача обновлена, id=%d", id)
}
