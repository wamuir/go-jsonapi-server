package handle

import (
	"context"
	"github.com/wamuir/go-jsonapi-server/graph"
	"github.com/wamuir/go-jsonapi-server/model"
	"net/http"
	"net/url"
)

// Resource is a handler for requests corresponding to a single resource
// (of type t string with identifier i string), with possible methods
// GET, HEAD, PATCH and DELETE.
func Resource(ctx context.Context, b url.URL, g graph.Graph, c model.Parameters, w http.ResponseWriter, r *http.Request, t, i string) (Response, *model.ErrorObject) {

	var response Response = NewResponse()

	q, errObj := model.ParseQueryString(r.URL, c)
	if errObj != nil {
		return response, errObj
	}

	switch r.Method {

	case "OPTIONS":

		_, errObj := model.GetResource(ctx, g, t, i, b, q)
		if errObj != nil {
			return response, errObj
		}
		response.Header.Set("Allow", "OPTIONS, GET, HEAD, PATCH, DELETE")
		response.Header.Set("Access-Control-Allow-Methods", "OPTIONS, GET, HEAD, PATCH, DELETE")
		response.Status = http.StatusNoContent
		return response, nil

	case "GET", "HEAD":

		document, errObj := model.GetResource(ctx, g, t, i, b, q)
		if errObj != nil {
			return response, errObj
		}
		response.Body = document
		response.Status = http.StatusOK
		return response, nil

	case "PATCH":

		// Method not implemented
		errObj := model.MakeError(http.StatusNotImplemented)
		errObj.Code = "b83e07"
		return response, errObj

	case "DELETE":

		errObj := model.DeleteResource(ctx, g, t, i)
		if errObj != nil {
			return response, errObj
		}
		response.Status = http.StatusNoContent
		return response, nil

	default:

		errObj := model.MakeError(http.StatusMethodNotAllowed)
		errObj.Code = "594414"
		return response, errObj

	}
}
