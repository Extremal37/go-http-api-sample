package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Extremal37/go-http-api-sample/internal/app/models"
	"io"
	"net/http"
)

type ContactJsonDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func contactToJSON(contact models.Contact) ContactJsonDTO {
	return ContactJsonDTO{
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
	}
}

func jsonToContact(json ContactJsonDTO) models.Contact {
	return models.Contact{
		FirstName: json.FirstName,
		LastName:  json.LastName,
	}
}

func (h *Handler) GetContacts(w http.ResponseWriter, r *http.Request) {
	contactsRaw, err := h.p.GetContacts(r.Context())
	if err != nil {
		h.WrapError(w, fmt.Errorf("unable to get contacts: %w", err))
		return
	}
	var contactsJSON []ContactJsonDTO

	for _, v := range contactsRaw {
		contactsJSON = append(contactsJSON, contactToJSON(v))
	}

	m := ResponseSuccess{
		Success: true,
		Result:  contactsJSON,
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

	var contactJSON ContactJsonDTO

	if err = json.Unmarshal(body, &contactJSON); err != nil {
		h.WrapBadRequest(w, fmt.Errorf("unable to parse contacat payload: %w", err))
		return
	}

	err = h.p.AddContact(r.Context(), jsonToContact(contactJSON))
	if err != nil {
		h.WrapError(w, fmt.Errorf("unable to add contact: %w", err))
		return
	}

	m := ResponseSuccess{
		Success: true,
		Result:  nil,
	}

	h.WrapNew(w, m)

}
