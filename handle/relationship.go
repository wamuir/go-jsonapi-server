package handle

import (
	"context"
	"github.com/wamuir/go-jsonapi-server/graph"
	"github.com/wamuir/go-jsonapi-server/model"
	"net/http"
	"net/url"
)

// Relationship is a handler for requests corresponding to a resource's
// relationship (resource of type t string and identifier i string with
// relationship keyed by k string), with possible methods GET, HEAD,
// PATCH and DELETE.
func Relationship(ctx context.Context, b url.URL, g graph.Graph, c model.Parameters, w http.ResponseWriter, r *http.Request, t, i, k string) (Response, *model.ErrorObject) {

	var response Response = NewResponse()

	q, errObj := model.ParseQueryString(r.URL, c)
	if errObj != nil {
		return response, errObj
	}

	switch r.Method {

	case "OPTIONS":

		_, errObj := model.GetRelationship(ctx, g, t, i, k, b, q)
		if errObj != nil {
			return response, errObj
		}
		response.Header.Set("Allow", "OPTIONS, GET, HEAD, POST, PATCH, DELETE")
		response.Header.Set("Access-Control-Allow-Methods", "OPTIONS, GET, HEAD, POST, PATCH, DELETE")
		response.Status = http.StatusNoContent
		return response, nil

	case "GET", "HEAD":

		document, errObj := model.GetRelationship(ctx, g, t, i, k, b, q)
		if errObj != nil {
			return response, errObj
		}
		response.Body = document
		response.Status = http.StatusOK
		return response, nil

	case "POST":

		// Validate content type
		errObj := validateMIME(r.Header.Get("Content-Type"))
		if errObj != nil {
			return response, errObj
		}

		// Parse request body
		document, errObj := model.Decode(r.Body)
		if errObj != nil {
			return response, errObj
		}

		// Post members to relationship
		errObj = model.PostRelationship(ctx, g, t, i, k, document)
		if errObj != nil {
			return response, errObj
		}

		// response.Header.Set("Location", i.URL(h))
		response.Status = http.StatusNoContent
		return response, nil

	case "PATCH":

		// Method not implemented
		errObj := model.MakeError(http.StatusNotImplemented)
		errObj.Code = "b4d563"
		return response, errObj

	case "DELETE":

		// Validate content type
		errObj := validateMIME(r.Header.Get("Content-Type"))
		if errObj != nil {
			return response, errObj
		}

		// Parse request body
		document, errObj := model.Decode(r.Body)
		if errObj != nil {
			return response, errObj
		}

		// Delete from relationship
		errObj = model.DeleteRelationship(ctx, g, t, i, k, document)
		if errObj != nil {
			return response, errObj
		}
		response.Status = http.StatusNoContent
		return response, nil

	default:

		// HTTP Method not allowed
		errObj := model.MakeError(http.StatusMethodNotAllowed)
		errObj.Code = "ce3a82"
		return response, errObj

	}
}
