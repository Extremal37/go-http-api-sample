package api

import (
	"github.com/Extremal37/go-http-api-sample/api/middleware"
	"github.com/Extremal37/go-http-api-sample/internals/app/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func CreateRoutes(handler *handlers.Handler, log *zap.SugaredLogger) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.GetRequestLogFunc(log))

	r.HandleFunc("/contact/list", handler.GetContacts).Methods("GET")
	r.HandleFunc("/contact", handler.AddContact).Methods("POST")

	r.MethodNotAllowedHandler = http.HandlerFunc(handler.WrapMethodNotAllowed)
	r.NotFoundHandler = http.HandlerFunc(handler.WrapNotFound)

	return r
}
