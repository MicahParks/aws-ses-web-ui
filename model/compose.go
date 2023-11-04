package model

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	jt "github.com/MicahParks/jsontype"
	"github.com/google/uuid"
)

type Compose struct {
	ID           int64                         `json:"-"`
	UUID         uuid.UUID                     `json:"uuid"`
	From         *jt.JSONType[*mail.Address]   `json:"from"`
	RecipientTo  []*jt.JSONType[*mail.Address] `json:"recipientTo"`
	RecipientCC  []*jt.JSONType[*mail.Address] `json:"recipientCC"`
	RecipientBCC []*jt.JSONType[*mail.Address] `json:"recipientBCC"`
	Subject      string                        `json:"subject"`
	Body         string                        `json:"body"`
	SESAccepted  bool                          `json:"sesAccepted"`
	Created      time.Time                     `json:"created"`
}

func (c Compose) DefaultsAndValidate() (Compose, error) {
	if c.From.Get() == nil {
		return c, fmt.Errorf("from is required: %w", jt.ErrDefaultsAndValidate)
	}
	if len(c.RecipientTo) == 0 && len(c.RecipientCC) == 0 && len(c.RecipientBCC) == 0 {
		return c, fmt.Errorf("at least one recipient is required: %w", jt.ErrDefaultsAndValidate)
	}
	normalizeEmails(c.RecipientTo)
	normalizeEmails(c.RecipientCC)
	normalizeEmails(c.RecipientBCC)
	c.Subject = strings.TrimSpace(c.Subject)
	if c.Subject == "" {
		return c, fmt.Errorf("subject is required: %w", jt.ErrDefaultsAndValidate)
	}
	if len(c.Subject) > 255 {
		return c, fmt.Errorf("subject must be less than 256 characters: %w", jt.ErrDefaultsAndValidate)
	}
	c.Body = strings.TrimSpace(c.Body)
	if c.Body == "" {
		return c, fmt.Errorf("body is required: %w", jt.ErrDefaultsAndValidate)
	}
	return c, nil
}

func normalizeEmails(arr []*jt.JSONType[*mail.Address]) {
	for i := range arr {
		arr[i] = jt.New[*mail.Address](must(mail.ParseAddress(strings.ToLower(strings.TrimSpace(arr[i].Get().Address)))))
	}
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
