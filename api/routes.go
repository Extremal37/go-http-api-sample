package api

import (
	"github.com/Extremal37/go-http-api-sample/api/middleware"
	"github.com/Extremal37/go-http-api-sample/internal/app/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type Handler interface {
	GetContacts(w http.ResponseWriter, r *http.Request)
	AddContact(w http.ResponseWriter, r *http.Request)
	WrapErrorWithStatus(w http.ResponseWriter, msg error, httpStatus int)
	WrapSuccessStatus(w http.ResponseWriter, m handlers.ResponseSuccess, httpStatus int)
	WrapNew(w http.ResponseWriter, m handlers.ResponseSuccess)
	WrapOK(w http.ResponseWriter, m handlers.ResponseSuccess)
	WrapBadRequest(w http.ResponseWriter, err error)
	WrapNotFound(w http.ResponseWriter, r *http.Request)
	WrapMethodNotAllowed(w http.ResponseWriter, r *http.Request)
}

func CreateRoutes(handler *handlers.Handler, log *zap.SugaredLogger) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.GetRequestLogFunc(log))

	r.HandleFunc("/contact/list", handler.GetContacts).Methods("GET")
	r.HandleFunc("/contact", handler.AddContact).Methods("POST")

	r.MethodNotAllowedHandler = http.HandlerFunc(handler.WrapMethodNotAllowed)
	r.NotFoundHandler = http.HandlerFunc(handler.WrapNotFound)

	return r
}
