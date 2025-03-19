package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Extremal37/go-http-api-sample/internal/app/models"
	"io"
	"net/http"
)

func (h *Handler) GetContacts(w http.ResponseWriter, r *http.Request) {
	contacts := h.p.GetContacts()

	m := Response{
		Success: true,
		Result:  contacts,
	}

	h.WrapOK(w, m)
}

func (h *Handler) AddContact(w http.ResponseWriter, r *http.Request) {
	// Read request body.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.WrapBadRequest(w, fmt.Errorf("unable to read request body: %w", err))
		return
	}

	var contact models.Contact

	if err = json.Unmarshal(body, &contact); err != nil {
		h.WrapBadRequest(w, fmt.Errorf("unable to parse contacat payload: %w", err))
		return
	}

	h.p.AddContact(contact)
	m := Response{
		Success: true,
		Result:  nil,
	}

	h.WrapNew(w, m)

}
