package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/wamuir/go-jsonapi-server/graph"
	"net/http"
	"net/url"
	"path"
)

type Resource struct {
	Identifier    string                 `json:"id,omitempty"`
	Type          string                 `json:"type"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
	Meta          map[string]interface{} `json:"meta,omitempty"`
	Relationships map[string]Document    `json:"relationships,omitempty"`
	// Links         map[string]string      `json:"links,omitempty"`
	Links LinksObject `json:"links,omitempty"`
}

func (resource Resource) Identify() Resource {

	return Resource{
		Type:       resource.Type,
		Identifier: resource.Identifier,
		Meta:       resource.Meta,
	}
}

func DeleteResource(ctx context.Context, g graph.Graph, t, i string) *ErrorObject {

	transaction, err := g.Transaction(ctx, false)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "fcd6d1"
		errObj.Title = "Encountered internal error while beginning graph transaction"
		errObj.Detail = err.Error()
		return errObj
	}
	defer transaction.Close()

	var tx *Tx = &Tx{transaction}

	errObj := tx.DeleteResource(t, i)
	if errObj != nil {
		return errObj
	}

	err = tx.Commit()
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "ff2a28"
		errObj.Title = "Encountered internal error while committing to graph"
		errObj.Detail = err.Error()
		return errObj
	}

	return nil
}

func (tx *Tx) DeleteResource(t, i string) *ErrorObject {

	err := tx.DeleteVertex(t, i)
	if err == graph.ErrNoRows {
		errObj := MakeError(http.StatusNotFound)
		errObj.Code = "eb476c"
		return errObj
	} else if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "14d479"
		errObj.Title = "Encounted internal error while deleting from graph"
		errObj.Detail = err.Error()
		return errObj
	}

	return nil
}

// Begin a new *Transaction and make a call to GetResourceLinkage()
func GetResource(ctx context.Context, g graph.Graph, t, i string, h url.URL, q QueryParams) (*Document, *ErrorObject) {

	var document *Document = &Document{}

	transaction, err := g.Transaction(ctx, true)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "6af933"
		errObj.Title = "Encountered internal error while beginning graph transaction"
		errObj.Detail = err.Error()
		return &Document{}, errObj
	}
	defer transaction.Close()

	var tx *Tx = &Tx{transaction}

	document, errObj := tx.GetResource(t, i, h, q)
	if errObj != nil {
		return document, errObj
	}

	return document, nil

}

func (tx *Tx) GetResource(t, i string, h url.URL, q QueryParams) (*Document, *ErrorObject) {

	document := &Document{}

	vertex, err := tx.FindVertex(t, i)
	if err != nil && err == graph.ErrNoRows {
		errObj := MakeError(http.StatusNotFound)
		errObj.Code = "bbf421"
		return document, errObj
	} else if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "dd6da6"
		errObj.Title = "Encountered internal error while querying graph"
		errObj.Detail = err.Error()
		return document, errObj
	}

	resource := Resource{
		Type:       vertex.Type,
		Identifier: vertex.Identifier,
	}

	err = json.Unmarshal(vertex.Attributes, &resource.Attributes)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "c40298"
		errObj.Title = "Encountered internal error while transforming data"
		errObj.Detail = err.Error()
		return document, errObj
	}

	err = json.Unmarshal(vertex.Meta, &resource.Meta)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "1da5e8"
		errObj.Title = "Encountered internal error while transforming data"
		errObj.Detail = err.Error()
		return document, errObj
	}

	edgeKeys, err := tx.FindDistinctEdgeKeys(t, i)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "443cda"
		errObj.Title = "Encountered internal error while querying graph"
		return document, errObj
	}

	/*
		includeKeys, keyRing, err := q.Include.SplitAndValidate(edgeKeys)
		if err != nil {
			return document, errors.New("Invalid `Include` entry")
		}

		q.Include = keyRing
	*/

	valid := q.Include.IsValidAgainst(edgeKeys)
	if !valid {
		errObj := MakeError(http.StatusBadRequest)
		errObj.Code = "9f9c9d"
		errObj.Title = "Invalid query string"
		errObj.Detail = "Unable to fulfill request for included resources"
		// This could hint include param if IsValidAgainst() returned err
		return document, errObj
	}

	relationships := make(map[string]Document)

	for _, k := range edgeKeys {

		var relationship Document

		if q.Include.Requests(k) { //stringInSlice(key, includeKeys) {

			relationship, errObj := tx.GetRelationship(t, i, k, h, q)
			if errObj != nil {
				return document, errObj
			}

			included := relationship.PopIncluded()

			relationships[k] = *relationship

			document.Included = document.Included.Merge(included)

		} else {

			/*
				relationship.Links = LinksObject{
					"self":    resource.URL(h) + "/relationships/" + k,
					"related": resource.URL(h) + "/" + k,
				}
				relationships[k] = relationship
			*/

			relationship.Links = LinksObject{}

			// Build "self" link
			ref, err := url.Parse(
				path.Join(resource.Type, resource.Identifier, "relationships", k),
			)
			if err != nil {
				errObj := MakeError(http.StatusInternalServerError)
				errObj.Code = "8a4660"
				errObj.Title = "Encountered internal error while generating response"
				errObj.Detail = err.Error()
				return document, errObj
			}

			relationship.Links["self"] = h.ResolveReference(ref).String()

			// Build "related" link
			ref, err = url.Parse(
				path.Join(resource.Type, resource.Identifier, k),
			)
			if err != nil {
				errObj := MakeError(http.StatusInternalServerError)
				errObj.Code = "0a4813"
				errObj.Title = "Encountered internal error while generating response"
				errObj.Detail = err.Error()
				return document, errObj
			}

			relationship.Links["related"] = h.ResolveReference(ref).String()

			relationships[k] = relationship

		}

	}

	if len(relationships) > 0 {
		resource.Relationships = relationships
	}

	ref, err := url.Parse(
		path.Join(resource.Type, resource.Identifier),
	)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "2aeacd"
		errObj.Title = "Encountered internal error while generating response"
		errObj.Detail = err.Error()
		return document, errObj
	}

	resource.Links = LinksObject{
		"self": h.ResolveReference(ref).String(),
	}

	document.Data = resource

	return document, nil

}

func PostResource(ctx context.Context, g graph.Graph, t string, d *Document) (Resource, *ErrorObject) {

	var identifier Resource

	transaction, err := g.Transaction(ctx, false)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "1c3686"
		errObj.Title = "Encountered internal error while beginning graph transaction"
		errObj.Detail = err.Error()
		return identifier, errObj
	}
	defer transaction.Close()

	var tx *Tx = &Tx{transaction}

	identifier, errObj := tx.PostResource(t, d)
	if errObj != nil {
		return identifier, errObj
	}

	err = tx.Commit()
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "bca977"
		errObj.Title = "Encountered internal error while committing graph transaction"
		errObj.Detail = err.Error()
		return identifier, errObj
	}

	return identifier, nil
}

func (tx *Tx) PostResource(t string, d *Document) (Resource, *ErrorObject) {

	var resource Resource

	resource, ok := d.Data.(Resource)
	if !ok {
		errObj := MakeError(http.StatusBadRequest)
		errObj.Code = "b7a83f"
		errObj.Title = "Bad request"
		errObj.Detail = fmt.Sprintf("Unable to assert data member as resource")
		return resource, errObj
	}

	if resource.Type != t {
		errObj := MakeError(http.StatusBadRequest)
		errObj.Code = "3b4ab2"
		errObj.Title = "Bad request"
		errObj.Detail = fmt.Sprintf("Resource of type %s cannot be posted to collection %s", resource.Type, t)
		return resource, errObj
	}

	if resource.Identifier == "" {
		resource.Identifier = uuid.New().String()
	}

	attributes, err := json.Marshal(resource.Attributes)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "830283"
		errObj.Title = "Encountered internal error while transforming data"
		errObj.Detail = err.Error()
		return resource, errObj
	}

	meta, err := json.Marshal(resource.Meta)
	if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "f2ae41"
		errObj.Title = "Encountered internal error while transforming data"
		errObj.Detail = err.Error()
		return resource, errObj
	}

	err = tx.InsertVertex(resource.Type, resource.Identifier, attributes, meta)
	if err == graph.ErrConflict {
		errObj := MakeError(http.StatusBadRequest)
		errObj.Code = "2910dd"
		errObj.Title = "Bad request"
		errObj.Detail = err.Error()
		return resource, errObj
	} else if err != nil {
		errObj := MakeError(http.StatusInternalServerError)
		errObj.Code = "bfb7ab"
		errObj.Title = "Encountered internal error while inserting data into graph"
		errObj.Detail = err.Error()
		return resource, errObj
	}

	for title, relationship := range resource.Relationships {

		errObj := tx.PostRelationship(resource.Type, resource.Identifier, title, &relationship)
		if errObj != nil {
			return resource, errObj
		}

	}

	return resource.Identify(), nil

}
