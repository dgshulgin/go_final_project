package server

import (
	"errors"
	"net/http"

	"github.com/dgshulgin/go_final_project/handler/server/dto"
	"github.com/dgshulgin/go_final_project/internal/nextdate"
)

// Обработчик GET "/api/nextdate?now=<20060102>&date=<20060102>&repeat=<правило>"
func (server TaskServer) NextDate(resp http.ResponseWriter, req *http.Request) {

	date := req.URL.Query().Get("date")
	now := req.URL.Query().Get("now")
	repeat := req.URL.Query().Get("repeat")

	rep := dto.RepeatCons{Date: date, Now: now, Repeat: repeat}
	server.logging(LogNextDate, rep)

	// Эта проверка добавлена из-за рассогласования в ТЗ:
	// endpoint /nextdate - рассматривает пустое поле repeat как ошибку
	// endpoint POST /api/task/ - в случае пустого repeat подставляет сегодняшее число.
	// Унифицировать проверку с помощью nextdate.validate невозможно.
	if len(repeat) == 0 {
		renderText(resp, http.StatusOK, "")
		return
	}

	err := nextdate.Validate(date, now, repeat)
	if err != nil {
		if !errors.Is(err, nextdate.ErrNextDateBeforeNow) {
			server.logging(err.Error(), nil)
			renderText(resp, http.StatusOK, "")
			return
		}
	}

	nextDate, err := nextdate.NextDate(date, now, repeat)
	if err != nil {
		msg := errors.Join(ErrNextDate, err).Error()
		server.logging(msg, nil)
		renderJSON(resp, http.StatusBadRequest, dto.Error{Error: msg})
		return
	}

	// Успех, возвращаем дату следующего события
	renderText(resp, http.StatusOK, nextDate)
}
