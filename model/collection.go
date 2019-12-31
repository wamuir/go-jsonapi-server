package model

import (
	"context"
	"fmt"
	"github.com/wamuir/go-jsonapi-server/graph"
	"net/http"
	"net/url"
)

type Collection []Resource

// Begin a new transaction (*Tx) and make a call to *Tx.GetCollection()
func GetCollection(ctx context.Context, g graph.Graph, t string, h url.URL, q QueryParams) (*Document, *ErrorObject) {

	var document *Document = &Document{}

	transaction, err := g.Transaction(ctx, true)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "001e49"
		errObj.Title = "Encountered internal error while beginning graph transaction"
		errObj.Detail = err.Error()
		return document, errObj
	}
	defer transaction.Close()

	var tx *Tx = &Tx{transaction}

	document, errObj := tx.GetCollection(t, h, q)
	if errObj != nil {
		return document, errObj
	}

	return document, nil
}

// Find and return a JSON:API document with a collection of resource objects as
// primary data.
func (tx *Tx) GetCollection(t string, h url.URL, q QueryParams) (*Document, *ErrorObject) {

	var document *Document = &Document{}

	count, err := tx.CountVertices(t)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "d71a15"
		errObj.Title = "Encountered internal error while querying graph"
		errObj.Detail = err.Error()
		return document, errObj
	}

	if count == 0 {
		errObj := MakeError(http.StatusNotFound)
		errObj.Code = "4aea6d"
		return document, errObj
	}

	collection := make(Collection, 0, count)

	vertices, err := tx.FindVertices(t, q.Limit, q.Offset)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "f3bce6"
		errObj.Title = "Encountered internal error while querying graph"
		errObj.Detail = err.Error()
		return document, errObj
	}

	for _, vertex := range vertices {

		identifier := Resource{
			Type:       vertex.Type,
			Identifier: vertex.Identifier,
		}

		resource, errObj := tx.GetResource(
			identifier.Type,
			identifier.Identifier,
			h,
			q,
		)
		if errObj != nil {
			return document, errObj
		}

		data, ok := resource.Data.(Resource)
		if !ok {
			errObj := MakeError(http.StatusInternalServerError)
			errObj.Code = "b5d595"
			errObj.Title = "Type assertion failed"
			errObj.Detail = fmt.Sprintf(
				"Interface is type %T not Resource",
				resource.Data,
			)
			return document, errObj
		}

		collection = append(collection, data)
		document.Included = document.Included.Merge(resource.Included)

	}

	document.Data = collection

	ref, err := url.Parse(t)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "2aeacd"
		errObj.Title = "Encountered internal error while generating response"
		errObj.Detail = err.Error()
		return document, errObj
	}

	document.Links = collection.paginate(
		h,
		ref,
		q.Limit,
		q.Offset,
		count,
	)

	return document, nil
}
