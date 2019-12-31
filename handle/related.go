package handle

import (
	"context"
	"github.com/wamuir/go-jsonapi-server/graph"
	"github.com/wamuir/go-jsonapi-server/model"
	"net/http"
	"net/url"
)

// Related is a handler for requests corresponding to a resource's
// related resources (primary resource of type t string and identifier
// i string with relationship keyed by k string), with possible
// methods GET and HEAD.
func Related(ctx context.Context, b url.URL, g graph.Graph, c model.Parameters, w http.ResponseWriter, r *http.Request, t, i, k string) (Response, *model.ErrorObject) {

	var response Response = NewResponse()

	q, errObj := model.ParseQueryString(r.URL, c)
	if errObj != nil {
		return response, errObj
	}

	switch r.Method {

	case "OPTIONS":

		_, errObj := model.GetRelated(ctx, g, t, i, k, b, q)
		if errObj != nil {
			return response, errObj
		}
		response.Header.Set("Allow", "OPTIONS, GET, HEAD")
		response.Header.Set("Access-Control-Allow-Methods", "OPTIONS, GET, HEAD")
		response.Status = http.StatusNoContent
		return response, nil

	case "GET", "HEAD":

		document, errObj := model.GetRelated(ctx, g, t, i, k, b, q)
		if errObj != nil {
			return response, errObj
		}
		response.Body = document
		response.Status = http.StatusOK
		return response, nil

	default:

		// HTTP Method not allowed
		errObj := model.MakeError(http.StatusMethodNotAllowed)
		errObj.Code = "ce3a82"
		return response, errObj

	}
}
