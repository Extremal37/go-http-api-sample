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
		errWrap := wrapError(w, fmt.Errorf("unable to get contacts: %w", err))
		if errWrap != nil {
			h.log.Warnf("unable to send error response for GetContacts func: %v", errWrap)
		}
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

	err = wrapOK(w, m)
	if err != nil {
		h.log.Warnf("unable to send success response for GetContacts func: %v", err)
	}
}

func (h *Handler) AddContact(w http.ResponseWriter, r *http.Request) {
	// Read request body.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errWrap := wrapBadRequest(w, fmt.Errorf("unable to read request body: %w", err))
		if errWrap != nil {
			h.log.Warnf("unable to send error response for AddContact func: %v", errWrap)
		}
		return
	}

	var contactJSON ContactJsonDTO

	if err = json.Unmarshal(body, &contactJSON); err != nil {
		errWrap := wrapBadRequest(w, fmt.Errorf("unable to parse contacat payload: %w", err))
		if errWrap != nil {
			h.log.Warnf("unable to send error response for AddContact func: %v", errWrap)
		}
		return
	}

	err = h.p.AddContact(r.Context(), jsonToContact(contactJSON))
	if err != nil {
		errWrap := wrapError(w, fmt.Errorf("unable to add contact: %w", err))
		if errWrap != nil {
			h.log.Warnf("unable to send error response for AddContact func: %v", errWrap)
		}
		return
	}

	m := ResponseSuccess{
		Success: true,
		Result:  nil,
	}

	err = wrapNew(w, m)
	if err != nil {
		h.log.Warnf("unable to send success response for GetContacts func: %v", err)
	}
}
