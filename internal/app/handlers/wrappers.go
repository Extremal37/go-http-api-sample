package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseWithError struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
}

func (h *Handler) WrapErrorWithStatus(w http.ResponseWriter, msg error, httpStatus int) {
	m := ResponseWithError{
		Success: false,
		Result:  msg.Error(),
	}

	res, err := json.Marshal(m)
	if err != nil {
		h.log.Warnf("Unable to encoding to JSON %v: %v", m, err)
	}
	// даем понять что ответ приходит в формате json
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(httpStatus)
	if _, err := fmt.Fprintln(w, string(res)); err != nil {
		h.log.Warnf("Unable to write HTTP response %d '%v' : %v", httpStatus, msg.Error(), err)
	}
}

type ResponseSuccess struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
}

// WrapSuccessStatus write map as json reply to httpWriter with success code.
func (h *Handler) WrapSuccessStatus(w http.ResponseWriter, m ResponseSuccess, httpStatus int) {
	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(httpStatus)

	if _, err := fmt.Fprintln(w, string(res)); err != nil {
		h.log.Warnf("Unable to write HTTP response %d : %v", httpStatus, err)
	}
}

func (h *Handler) WrapNew(w http.ResponseWriter, m ResponseSuccess) {
	h.WrapSuccessStatus(w, m, http.StatusCreated)
}

func (h *Handler) WrapOK(w http.ResponseWriter, m ResponseSuccess) {
	h.WrapSuccessStatus(w, m, http.StatusOK)
}

// WrapBadRequest write json error to httpWriter handler with 400 code.
func (h *Handler) WrapBadRequest(w http.ResponseWriter, err error) {
	h.WrapErrorWithStatus(w, err, http.StatusBadRequest)
}

// WrapNotFound write 'not found' as json to httpWriter with 404 code.
func (h *Handler) WrapNotFound(w http.ResponseWriter, r *http.Request) {
	// Не печатай это, пожалуйста, потому что у нас 404 используется в случае, если объект пустой, например или ответ пустой (нет плана проверки).
	//h.log.Debugf("Incorrect URL: %v", r.URL)
	h.WrapErrorWithStatus(w, fmt.Errorf("not found"), http.StatusNotFound)
}

// WrapMethodNotAllowed write 'method not allowed' as json to httpWriter with 405 code.
func (h *Handler) WrapMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	h.log.Debugf("Incorrect method: %v", r.Method)
	h.WrapErrorWithStatus(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
}

// WrapError write json error to httpWriter handler with 500 code.
func (h *Handler) WrapError(w http.ResponseWriter, err error) {
	h.log.Warnf("Internal error: %v", err)
	h.WrapErrorWithStatus(w, err, http.StatusInternalServerError)
}
