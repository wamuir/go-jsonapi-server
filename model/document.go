package model

import (
	"bytes"
	"encoding/json"
	"github.com/wamuir/go-jsonapi-server/graph"
)

const jsonapi_version = "1.0" // Version of JSON:API implementation

type Tx struct{ graph.Tx }

type Edge = graph.Edge

type Vertex = graph.Vertex

type Document struct {
	JSONAPI  *Implementation        `json:"jsonapi,omitempty"`
	Data     interface{}            `json:"data,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
	Links    LinksObject            `json:"links,omitempty"`
	Included Included               `json:"included,omitempty"`
	Errors   []ErrorObject          `json:"errors,omitempty"`
}

type Implementation struct {
	Version string                 `json:"version,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

type LinksObject map[string]interface{}

type Link struct {
	Href string                 `json:"href,omitempty"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}

func NewDocument() Document {
	var document Document
	document.Version()
	return document
}

func (document *Document) ContentLength() int {

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	encoder.Encode(document)

	return buf.Len()
}

func (document *Document) PopIncluded() Included {

	var included Included = document.Included
	document.Included = Included{}

	return included
}

func (document *Document) Version() {

	document.JSONAPI = &Implementation{
		Version: jsonapi_version,
		Meta: map[string]interface{}{
			"links": map[string]interface{}{
				"self": Link{
					Href: "http://jsonapi.org/format/1.0/",
				},
			},
		},
	}

	return
}
