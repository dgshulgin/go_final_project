package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	webDir    = "./web"
	rootRoute = "/"
)

func Router(log logrus.FieldLogger) (*mux.Router, error) {
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(webDir)))
	return router, nil
}
