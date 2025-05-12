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

func wrapErrorWithStatus(w http.ResponseWriter, msg error, httpStatus int) error {
	m := ResponseWithError{
		Success: false,
		Result:  msg.Error(),
	}

	res, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("unable to encoding to JSON %v: %v", m, err)
	}
	// даем понять что ответ приходит в формате json
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(httpStatus)
	if _, err := fmt.Fprintln(w, string(res)); err != nil {

		return fmt.Errorf("unable to write HTTP response %d '%v' : %v", httpStatus, msg.Error(), err)
	}
	return nil
}

type ResponseSuccess struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
}

// WrapSuccessStatus write map as json reply to httpWriter with success code.
func wrapSuccessStatus(w http.ResponseWriter, m ResponseSuccess, httpStatus int) error {
	res, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("unable to encoding to JSON %v: %v", m, err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(httpStatus)

	if _, err := fmt.Fprintln(w, string(res)); err != nil {
		return fmt.Errorf("unable to write HTTP response %d : %v", httpStatus, err)
	}
	return nil
}

func wrapNew(w http.ResponseWriter, m ResponseSuccess) error {
	err := wrapSuccessStatus(w, m, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}
	return nil
}

func wrapOK(w http.ResponseWriter, m ResponseSuccess) error {
	err := wrapSuccessStatus(w, m, http.StatusOK)
	if err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}
	return nil
}

// WrapBadRequest write json error to httpWriter handler with 400 code.
func wrapBadRequest(w http.ResponseWriter, sendErr error) error {
	err := wrapErrorWithStatus(w, sendErr, http.StatusBadRequest)
	if err != nil {
		return fmt.Errorf("failed to send BadRequest: %w", err)
	}
	return nil
}

// WrapNotFound write 'not found' as json to httpWriter with 404 code.
func (h *Handler) WrapNotFound(w http.ResponseWriter, r *http.Request) {
	_ = r
	err := wrapErrorWithStatus(w, fmt.Errorf("not found"), http.StatusNotFound)
	if err != nil {
		h.log.Warnf("failed to send NotFound: %v", err)
	}
}

// WrapMethodNotAllowed write 'method not allowed' as json to httpWriter with 405 code.
func (h *Handler) WrapMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	h.log.Debugf("Incorrect method: %v", r.Method)
	err := wrapErrorWithStatus(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
	if err != nil {
		h.log.Warnf("failed to send MethodNotAllowed: %v", err)
	}
}

// WrapError write json error to httpWriter handler with 500 code.
func wrapError(w http.ResponseWriter, errSend error) error {
	err := wrapErrorWithStatus(w, errSend, http.StatusInternalServerError)
	if err != nil {
		return fmt.Errorf("internal error: %v", err)
	}
	return nil
}
