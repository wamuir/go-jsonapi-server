package model

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/wamuir/go-jsonapi-core"
	"github.com/wamuir/go-jsonapi-server/graph"
)

func DeleteRelationship(ctx context.Context, g graph.Graph, t, i, k string, d *core.Document) *core.Error {

	transaction, err := g.Transaction(ctx, true)
	if err != nil {
		errObj := core.MakeError(http.StatusInternalServerError)
		errObj.Code = "8c58b4"
		errObj.Title = "Encountered internal error while beginning graph transaction"
		errObj.Detail = err.Error()
		return errObj
	}
	defer transaction.Close()

	var tx *Tx = &Tx{transaction}

	errObj := tx.DeleteRelationship(t, i, k, d)
	if errObj != nil {
		return errObj
	}

	return nil
}

func (tx *Tx) DeleteRelationship(t, i, k string, document *core.Document) *core.Error {

	var collection []core.Resource

	switch v := document.Data.(type) {

	case core.Collection:
		collection = v

	case core.Resource:
		collection = []core.Resource{v}

	default:
		errObj := core.MakeError(http.StatusBadRequest)
		errObj.Code = "207867"
		errObj.Title = "Bad Request"
		errObj.Detail = "Unable to assert data member as collection or resource"
		return errObj

	}

	for _, resource := range collection {

		err := tx.DeleteEdge(t, i, resource.Type, resource.Identifier, k)
		if err == graph.ErrNoRows {
			// pass, per JSON:API spec
		} else if err != nil {
			errObj := core.MakeError(http.StatusInternalServerError)
			errObj.Code = "c0e905"
			errObj.Title = "Encounted internal error while deleting from graph"
			errObj.Detail = err.Error()
			return errObj
		}

	}

	return nil
}

func GetRelationship(ctx context.Context, g graph.Graph, t, i, k string, h url.URL, q QueryParams) (*core.Document, *core.Error) {

	var document *core.Document = &core.Document{}

	transaction, err := g.Transaction(ctx, true)
	if err != nil {
		errObj := core.MakeError(http.StatusInternalServerError)
		errObj.Code = "56d959"
		errObj.Title = "Encountered internal error while beginning graph transaction"
		errObj.Detail = err.Error()
		return document, errObj
	}
	defer transaction.Close()

	var tx *Tx = &Tx{transaction}

	document, errObj := tx.GetRelationship(t, i, k, h, q)
	if errObj != nil {
		return document, errObj
	}

	return document, nil
}

func (tx *Tx) GetRelationship(t, i, k string, h url.URL, q QueryParams) (*core.Document, *core.Error) {

	var document *core.Document = &core.Document{}

	// Count of related vertices is needed for pagination
	count, err := tx.CountRelatedVertices(t, i, k)
	if err != nil {
		errObj := core.MakeError(http.StatusInternalServerError)
		errObj.Code = "f3bc34"
		errObj.Title = "Encountered internal error while querying graph"
		errObj.Detail = err.Error()
		return document, errObj
	}

	collection := make(core.Collection, 0, count)

	var edges []Edge
	edges, err = tx.FindEdges(t, i, k, q.Limit, q.Offset)
	if err != nil {
		errObj := core.MakeError(http.StatusInternalServerError)
		errObj.Code = "38440d"
		errObj.Title = "Encountered internal error while querying graph"
		errObj.Detail = err.Error()
		return document, errObj
	}

	for _, edge := range edges {

		resource := core.Resource{
			Type:       edge.To.Type,
			Identifier: edge.To.Identifier,
		}

		err := json.Unmarshal(edge.Meta, &resource.Meta)
		if err != nil {
			errObj := core.MakeError(http.StatusInternalServerError)
			errObj.Code = "286c1b"
			errObj.Title = "Encountered internal error while transforming data"
			errObj.Detail = err.Error()
			return document, errObj
		}

		collection = append(collection, resource)

	}

	switch {

	case count == 0:

		errObj := core.MakeError(http.StatusNotFound)
		errObj.Code = "9fbdd5"
		return document, errObj

	case count == 1:

		identifier := collection[0]

		if q.Include.Requests(k) {
			q.Include = q.Include.SplitOn(k)
			resource, errObj := tx.GetResource(
				identifier.Type,
				identifier.Identifier,
				h,
				q,
			)
			if errObj != nil {
				return document, errObj
			}

			data, ok := resource.Data.(core.Resource)
			if !ok {
				errObj := core.MakeError(http.StatusInternalServerError)
				errObj.Code = "f8f4c4"
				errObj.Title = "Type assertion failed"
				errObj.Detail = fmt.Sprintf(
					"Interface of type %T is not Resource",
					resource.Data,
				)
				return document, errObj
			}

			document.Included = document.Included.MergeResource(data)
			document.Included = document.Included.Merge(resource.Included)
		}

		document.Data = identifier

	case count >= 2:

		if q.Include.Requests(k) {
			q.Include = q.Include.SplitOn(k)
			for _, identifier := range collection {
				resource, errObj := tx.GetResource(
					identifier.Type,
					identifier.Identifier,
					h,
					q,
				)
				if errObj != nil {
					return document, errObj
				}

				data, ok := resource.Data.(core.Resource)
				if !ok {
					errObj := core.MakeError(http.StatusInternalServerError)
					errObj.Code = "f7c124"
					errObj.Title = "Type assertion failed"
					errObj.Detail = fmt.Sprintf(
						"Interface of type %T is not Resource",
						resource.Data,
					)
					return document, errObj
				}

				document.Included = document.Included.MergeResource(data)
				document.Included = document.Included.Merge(resource.Included)
			}
		}

		document.Data = collection

		ref, err := url.Parse(
			path.Join(t, i, "relationships", k),
		)
		if err != nil {
			errObj := core.MakeError(http.StatusInternalServerError)
			errObj.Code = "1d385a"
			errObj.Title = "Encountered internal error while generating response"
			errObj.Detail = err.Error()
			return document, errObj
		}

		_, document.Links = paginate(
			collection,
			h,
			ref,
			q.Limit,
			q.Offset,
			count,
		)

	}

	return document, nil
}

func PostRelationship(ctx context.Context, g graph.Graph, t, i, k string, document *core.Document) *core.Error {

	transaction, err := g.Transaction(ctx, false)
	if err != nil {
		errObj := core.MakeError(http.StatusInternalServerError)
		errObj.Code = "5ea87b"
		errObj.Title = "Encountered internal error while beginning graph transaction"
		errObj.Detail = err.Error()
		return errObj
	}
	defer transaction.Close()

	var tx *Tx = &Tx{transaction}

	errObj := tx.PostRelationship(t, i, k, document)
	if errObj != nil {
		return errObj
	}

	err = tx.Commit()
	if err != nil {
		errObj := core.MakeError(http.StatusInternalServerError)
		errObj.Code = "6f6ee0"
		errObj.Title = "Encountered internal error while committing to graph"
		errObj.Detail = err.Error()
		return errObj
	}

	return nil
}

func (tx *Tx) PostRelationship(t, i, k string, document *core.Document) *core.Error {

	var collection []core.Resource

	m, err := decodeDataMbr(document.Data)
	if err != nil {
		return err
	}

	switch v := m.(type) {
	// switch v := document.Data.(type) {

	case core.Collection:
		collection = v

	case core.Resource:
		collection = []core.Resource{v}

	default:
		errObj := core.MakeError(http.StatusBadRequest)
		errObj.Code = "cebc7c"
		errObj.Title = "Bad Request"
		errObj.Detail = "Unable to assert data member as collection or resource"
		return errObj

	}

	for pos, related := range collection {

		// Marshal meta member
		meta, err := json.Marshal(related.Meta)
		if err != nil {
			errObj := core.MakeError(http.StatusInternalServerError)
			errObj.Code = "b1474d"
			errObj.Title = "Encountered internal error while transforming data"
			errObj.Detail = err.Error()
			return errObj
		}

		err = tx.InsertEdge(
			t,
			i,
			related.Type,
			related.Identifier,
			k,
			pos,
			meta,
		)
		if err == graph.ErrConflict {
			// Pass if resource is already in relationship
		} else if err == graph.ErrNoRows {
			errObj := core.MakeError(http.StatusNotFound)
			errObj.Code = "54132b"
			return errObj
		} else if err != nil {
			errObj := core.MakeError(http.StatusInternalServerError)
			errObj.Code = "c2589a"
			errObj.Title = "Encountered internal error while inserting into graph"
			errObj.Detail = err.Error()
			return errObj
		}
	}

	return nil
}
