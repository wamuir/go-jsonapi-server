package handle

import (
	"context"
	"github.com/wamuir/go-jsonapi-server/graph"
	"github.com/wamuir/go-jsonapi-server/model"
	"net/http"
	"net/url"
	"path"
)

// Collection is a handler for requests corresponding to a collection of
// resources (of type t string), with possible methods GET, HEAD and POST.
func Collection(ctx context.Context, b url.URL, g graph.Graph, c model.Parameters, w http.ResponseWriter, r *http.Request, t string) (Response, *model.ErrorObject) {

	var response Response = NewResponse()

	q, errObj := model.ParseQueryString(r.URL, c)
	if errObj != nil {
		return response, errObj
	}

	switch r.Method {

	case "OPTIONS":

		// Verify collection exists
		_, errObj := model.GetCollection(ctx, g, t, b, q)
		if errObj != nil {
			return response, errObj
		}
		response.Header.Set("Allow", "OPTIONS, GET, HEAD, POST")
		response.Header.Set("Access-Control-Allow-Methods", "OPTIONS, GET, HEAD, POST")
		response.Status = http.StatusNoContent
		return response, nil

	case "GET", "HEAD":

		// Get collection
		document, errObj := model.GetCollection(ctx, g, t, b, q)
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

		// Post new resource
		i, errObj := model.PostResource(ctx, g, t, document)
		if errObj != nil {
			return response, errObj
		}

		// Get the resource
		document, errObj = model.GetResource(ctx, g, i.Type, i.Identifier, b, q)
		if errObj != nil {
			return response, errObj
		}

		// Build link for the new resource
                ref, err := url.Parse(
		       path.Join(i.Type, i.Identifier),
	        )
		if err != nil {
		      errObj := model.MakeError(http.StatusInternalServerError)
		      errObj.Code = "bfe23f"
		      errObj.Title = "Encountered internal error while generating response"
		      errObj.Detail = err.Error()
	        }

		response.Body = document
		response.Header.Set("Location", b.ResolveReference(ref).String())
		response.Status = http.StatusCreated
		return response, nil

	default:

		// HTTP Method not allowed
		errObj := model.MakeError(http.StatusMethodNotAllowed)
		errObj.Code = "8e5fce"
		return response, errObj

	}
}
