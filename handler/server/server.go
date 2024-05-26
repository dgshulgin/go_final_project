package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dgshulgin/go_final_project/internal/repository"
	"github.com/sirupsen/logrus"
)

var (
	ErrBadRequest        = errors.New("ошибка чтения данных запроса")
	ErrRequestValidation = errors.New("ошибка валидации данных")
	ErrSaveData          = errors.New("ошибка сохранения данных")
	ErrEmptyTitle        = errors.New("поле title не может быть пустым")
	ErrNoId              = errors.New("не указан идентификатор задачи")
	ErrInvalidId         = errors.New("некорректный идентификатор задачи")
	ErrNotExistId        = errors.New("идентификатор не существует")
	ErrDelete            = errors.New("во время удаления задачи возникла ошибка")
	ErrSelect            = errors.New("ошибка при запросе к БД")
	ErrNextDate          = errors.New("ошибка при вычислении даты события")
	ErrSearch            = errors.New("ошибка во время поиска данных")
	ErrInvalidJWToken    = errors.New("некорректный токен")
	ErrAuthentication    = errors.New("ошибка аутентификации пользователя")
	ErrWrongPassword     = errors.New("неверный пароль")
	ErrCreateJWT         = errors.New("ошибка при формировании JWT")
	ErrTokenVerification = errors.New("ошибка валидации JWT-токена")
)

var (
	LogCreateTask     = string("создание задачи с параметрами %v")
	LogDeleteTaskById = string("удаление задачи id=%d")
	LogUpdateTask     = string("обновление задачи с параметрами %v")
	LogCompleteTask   = string("запрос на завершение задачи id=%d")
	LogGetTaskById    = string("запрос информации о задаче id=%d")
	LogNextDate       = string("вычисление следующей даты события, условия %v")
	LogSearchText     = string("поиск по тексту заголовка или комментария \"%s\"")
	LogSearchDate     = string("поиск по дате %s")
)

type TaskServer struct {
	log  logrus.FieldLogger
	repo *repository.Repository
}

func NewTaskServer(log logrus.FieldLogger, repo *repository.Repository) *TaskServer {
	return &TaskServer{log, repo}
}

func (ts TaskServer) logging(msg string, v interface{}) {
	if v != nil {
		data, _ := json.Marshal(v)
		ts.log.Debugf(msg, string(data))
		return
	}
	ts.log.Debugf(msg)
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
