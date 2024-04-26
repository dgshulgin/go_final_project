package handler

import (
	"net/http"

	"github.com/dgshulgin/go_final_project/internal/pkg/nextdate"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	webDir    = "./web"
	rootRoute = "/"
)

func Router(log logrus.FieldLogger) (*mux.Router, error) {
	router := mux.NewRouter()

	router.HandleFunc("/api/nextdate", func(resp http.ResponseWriter, req *http.Request) {
		start := req.FormValue("date")
		now := req.FormValue("now")
		repeat := req.FormValue("repeat")
		log.Printf("api/nextdate, start=%s, now=%s, repeat={%s}", start, now, repeat)
		if len(repeat) == 0 {
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte{})
			log.Printf("отправлена пустая строка") // по логике, тако таск должен удаляться
			return
		}
		nextDate, err := nextdate.NextDate(start, now, repeat)
		if err != nil {
			log.Printf("ошибка, %q", err)
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write([]byte(err.Error()))
			return
		}
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(nextDate))
		log.Printf("отправлено, %q", nextDate)
	}).Methods(http.MethodGet)

	// я должен быть последним
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(webDir)))

	return router, nil
}
