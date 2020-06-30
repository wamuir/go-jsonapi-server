package model

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/wamuir/go-jsonapi-core"
	"github.com/wamuir/go-jsonapi-server/graph"
)

// Begin a new transaction (*Tx) and make a call to *Tx.GetCollection()
func GetCollection(ctx context.Context, g graph.Graph, t string, h url.URL, q QueryParams) (*core.Document, *core.Error) {

	transaction, err := g.Transaction(ctx, true)
	if err != nil {
		e := core.MakeError(http.StatusInternalServerError)
		e.Code = "001e49"
		e.Title = "Encountered internal error while beginning graph transaction"
		e.Detail = err.Error()
		return nil, e
	}
	defer transaction.Close()

	var tx *Tx = &Tx{transaction}

	document, errObj := tx.GetCollection(t, h, q)
	if errObj != nil {
		return nil, errObj
	}

	return document, nil
}

// Find and return a JSON:API document with a collection of resource objects as
// primary data.
func (tx *Tx) GetCollection(t string, h url.URL, q QueryParams) (*core.Document, *core.Error) {

	var document *core.Document = &core.Document{}

	count, err := tx.CountVertices(t)
	if err != nil {
		e := core.MakeError(http.StatusInternalServerError)
		e.Code = "d71a15"
		e.Title = "Encountered internal error while querying graph"
		e.Detail = err.Error()
		return nil, e
	}

	if count == 0 {
		e := core.MakeError(http.StatusNotFound)
		e.Code = "4aea6d"
		return nil, e
	}

	collection := make(core.Collection, 0, count)

	vertices, err := tx.FindVertices(t, q.Limit, q.Offset)
	if err != nil {
		e := core.MakeError(http.StatusInternalServerError)
		e.Code = "f3bce6"
		e.Title = "Encountered internal error while querying graph"
		e.Detail = err.Error()
		return nil, e
	}

	for _, vertex := range vertices {

		identifier := core.Resource{
			Type:       vertex.Type,
			Identifier: vertex.Identifier,
		}

		resource, modelErr := tx.GetResource(
			identifier.Type,
			identifier.Identifier,
			h,
			q,
		)
		if modelErr != nil {
			return nil, modelErr
		}

		data, ok := resource.Data.(core.Resource)
		if !ok {
			e := core.MakeError(http.StatusInternalServerError)
			e.Code = "b5d595"
			e.Title = "Type assertion failed"
			e.Detail = fmt.Sprintf(
				"Interface is type %T not Resource",
				resource.Data,
			)
			return nil, e
		}

		collection = append(collection, data)
		document.Included = document.Included.Merge(resource.Included)

	}

	document.Data = collection

	ref, err := url.Parse(t)
	if err != nil {
		e := core.MakeError(http.StatusInternalServerError)
		e.Code = "2aeacd"
		e.Title = "Encountered internal error while generating response"
		e.Detail = err.Error()
		return nil, e
	}

	_, document.Links = paginate(
		collection,
		h,
		ref,
		q.Limit,
		q.Offset,
		count,
	)

	return document, nil
}
