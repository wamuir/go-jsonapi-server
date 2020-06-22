package model

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/wamuir/go-jsonapi-server/graph"
)

func GetRelated(ctx context.Context, g graph.Graph, t, i, k string, h url.URL, q QueryParams) (*Document, *ModelError) {

	var document *Document = &Document{}

	transaction, err := g.Transaction(ctx, true)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "13f87d"
		errObj.Title = "Encountered internal error while beginning graph transaction"
		errObj.Detail = err.Error()
		return &Document{}, errObj
	}
	defer transaction.Close()

	var tx *Tx = &Tx{transaction}

	document, errObj := tx.GetRelated(t, i, k, h, q)
	if errObj != nil {
		return document, errObj
	}

	return document, nil
}

func (tx *Tx) GetRelated(t, i, k string, h url.URL, q QueryParams) (*Document, *ModelError) {

	var document *Document = &Document{}

	count, err := tx.CountRelatedVertices(t, i, k)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "e9bf2d"
		errObj.Title = "Encountered internal error while querying graph"
		errObj.Detail = err.Error()
		return document, errObj
	}

	collection := make(Collection, 0, count)

	if count == 0 {
		errObj := MakeError(http.StatusNotFound)
		errObj.Code = "9fbdd5"
		return document, errObj
	}

	var edges []Edge
	edges, err = tx.FindEdges(t, i, k, q.Limit, q.Offset)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "9a8ffa"
		errObj.Title = "Encountered internal error while querying graph"
		errObj.Detail = err.Error()
		return document, errObj
	}

	for _, edge := range edges {

		identifier := Resource{
			Type:       edge.To.Type,
			Identifier: edge.To.Identifier,
		}

		err := json.Unmarshal(edge.Meta, &identifier.Meta)
		if err != nil {
			errObj := MakeError(http.StatusInternalServerError)
			errObj.Code = "67590d"
			errObj.Title = "Encountered internal error while transforming data"
			errObj.Detail = err.Error()
			return document, errObj
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
			errObj.Code = "9cf16c"
			errObj.Title = "Type assertion failed"
			errObj.Detail = fmt.Sprintf(
				"Interface of type %T is not Resource",
				resource.Data,
			)
			return document, errObj
		}

		data.Meta = identifier.Meta

		collection = append(collection, data)

		document.Included = document.Included.Merge(resource.Included)

	}

	document.Data = collection

	ref, err := url.Parse(
		path.Join(t, i, k),
	)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "1d385a"
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
