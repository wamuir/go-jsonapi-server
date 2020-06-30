package handle

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/wamuir/go-jsonapi-core"
	"github.com/wamuir/go-jsonapi-server/model"
)

// Relationship is a handler for requests corresponding to a resource's
// relationship (resource of type t string and identifier i string with
// relationship keyed by k string), with possible methods GET, HEAD,
// PATCH and DELETE.
func (env *Environment) HandleRelationship(w http.ResponseWriter, r *http.Request) {

	var response Response = NewResponse()

	q, e := model.ParseQueryString(r.URL, env.Parameters)
	if e != nil {
		env.Fail(w, r, e)
		return
	}

	t := chi.URLParam(r, "type")
	i := chi.URLParam(r, "id")
	k := chi.URLParam(r, "relationship")

	switch r.Method {

	case "OPTIONS":

		_, e := model.GetRelationship(r.Context(), env.Graph, t, i, k, env.BaseURL, q)
		if e != nil {
			env.Fail(w, r, e)
			return
		}
		response.Header.Set("Allow", "OPTIONS, GET, HEAD, POST, PATCH, DELETE")
		response.Header.Set("Access-Control-Allow-Methods", "OPTIONS, GET, HEAD, POST, PATCH, DELETE")
		response.Status = http.StatusNoContent
		env.Success(w, r, response)
		return

	case "GET", "HEAD":

		document, e := model.GetRelationship(r.Context(), env.Graph, t, i, k, env.BaseURL, q)
		if e != nil {
			env.Fail(w, r, e)
			return
		}
		response.Body = document
		response.Status = http.StatusOK
		env.Success(w, r, response)
		return

	case "POST":

		// Validate content type
		e := validateMIME(r.Header.Get("Content-Type"))
		if e != nil {
			env.Fail(w, r, e)
			return
		}

		// Parse request body
		document, e := model.Decode(r.Body)
		if e != nil {
			env.Fail(w, r, e)
			return
		}

		// Post members to relationship
		e = model.PostRelationship(r.Context(), env.Graph, t, i, k, document)
		if e != nil {
			env.Fail(w, r, e)
			return
		}

		// response.Header.Set("Location", i.URL(h))
		response.Status = http.StatusNoContent
		env.Success(w, r, response)
		return

	case "PATCH":

		// Method not implemented
		e := core.MakeError(http.StatusNotImplemented)
		e.Code = "b4d563"
		env.Fail(w, r, e)
		return

	case "DELETE":

		// Validate content type
		e := validateMIME(r.Header.Get("Content-Type"))
		if e != nil {
			env.Fail(w, r, e)
			return
		}

		// Parse request body
		document, e := model.Decode(r.Body)
		if e != nil {
			env.Fail(w, r, e)
			return
		}

		// Delete from relationship
		e = model.DeleteRelationship(r.Context(), env.Graph, t, i, k, document)
		if e != nil {
			env.Fail(w, r, e)
			return
		}
		response.Status = http.StatusNoContent
		env.Success(w, r, response)
		return

	default:

		// HTTP Method not allowed
		e := core.MakeError(http.StatusMethodNotAllowed)
		e.Code = "ce3a82"
		env.Fail(w, r, e)
		return

	}
}
