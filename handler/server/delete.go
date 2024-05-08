package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
)

// Обработчик DELETE /api/task/done?id=<идентификатор>
func (server TaskServer) Delete(resp http.ResponseWriter, req *http.Request) {

	id0 := req.URL.Query().Get("id")
	if len(id0) == 0 {
		msg := "не указан идентификатор"
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	server.log.Printf("запрос на удаление задачи id=%s", id0)

	id, err := strconv.Atoi(id0)
	if err != nil {
		msg := fmt.Sprintf("некорректный идентификатор %d", id)
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	err = server.repo.Delete([]uint{uint(id)})
	if err != nil {
		msg := fmt.Sprintf("ошибка при удалении задачи id=%d", id)
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	renderJSON(resp, http.StatusOK, dto.Ok{})
	server.log.Infof("задача удалена id=%d", id)
}
