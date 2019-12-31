package model

import (
	"net/http"
	"strconv"
)

type ErrorObject struct {
	Identifier string                 `json:"id,omitempty"`
	Links      *LinksObject           `json:"links,omitempty"`
	Status     string                 `json:"status,omitempty"`
	Code       string                 `json:"code,omitempty"`
	Title      string                 `json:"title,omitempty"`
	Detail     string                 `json:"detail,omitempty"`
	Source     *SourceObject          `json:"source,omitempty"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
}

type SourceObject struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

func MakeError(status int) *ErrorObject {

	return &ErrorObject{
		Status: strconv.Itoa(status),
		Title:  http.StatusText(status),
	}
}
