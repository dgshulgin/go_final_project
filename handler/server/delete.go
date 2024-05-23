package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
)

// Обработчик DELETE /api/task/done?id=<идентификатор>
func (server TaskServer) Delete(resp http.ResponseWriter, req *http.Request) {

	id0 := req.URL.Query().Get("id")
	if len(id0) == 0 {
		server.logging(ErrNoId.Error(), nil)
		renderJSON(resp,
			http.StatusBadRequest, dto.Error{Error: ErrNoId.Error()})
		return
	}

	id, err := strconv.Atoi(id0)
	if err != nil {
		server.logging(ErrInvalidId.Error(), nil)
		renderJSON(resp,
			http.StatusBadRequest, dto.Error{Error: ErrInvalidId.Error()})
		return
	}

	server.logging(LogDeleteTaskById, id)

	err = server.repo.Delete([]uint{uint(id)})
	if err != nil {
		msg := errors.Join(ErrDelete, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	// Успех, возвращается пустой JSON
	renderJSON(resp, http.StatusOK, dto.Ok{})
}
