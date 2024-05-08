package server

import (
	"fmt"
	"net/http"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
	"github.com/dgshulgin/go_final_project/internal/nextdate"
)

// Обработчик GET "/api/nextdate?now=<20060102>&date=<20060102>&repeat=<правило>"
func (server TaskServer) NextDate(resp http.ResponseWriter, req *http.Request) {

	date := req.URL.Query().Get("date")
	now := req.URL.Query().Get("now")
	repeat := req.URL.Query().Get("repeat")

	server.log.Printf("проверка функции nextdate, данные {date=%s, now=%s, repeat=%s}",
		date, now, repeat)

	if len(repeat) == 0 {
		server.log.Printf("поле repeat не определено, возвращается пустая строка")
		renderText(resp, http.StatusOK, "")
		return
	}

	nextDate, err := nextdate.NextDate(date, now, repeat)
	if err != nil {
		msg := fmt.Sprintf("невозможно вычислить дату, %s", err.Error())
		server.log.Errorf(msg)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	renderText(resp, http.StatusOK, nextDate)
	server.log.Printf("новая дата {date=%s}", nextDate)
}
