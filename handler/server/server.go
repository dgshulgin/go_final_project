package server

import (
	"encoding/json"
	"net/http"

	"github.com/dgshulgin/go_final_project/internal/repository"
	"github.com/sirupsen/logrus"
)

type TaskServer struct {
	log  logrus.FieldLogger
	repo *repository.Repository
}

func NewTaskServer(log logrus.FieldLogger, repo *repository.Repository) *TaskServer {
	return &TaskServer{log, repo}
}

// Обертка, для отправки клиенту сообщения в формате JSON
func renderJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)

	return err
}

// Обертка, для отправки клиенту сообщения в формате plain text
// Используется только в GET /api/nextdate/
func renderText(w http.ResponseWriter, status int, v string) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, err := w.Write([]byte(v))
	return err
}
