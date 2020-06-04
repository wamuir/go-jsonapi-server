package handle

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/wamuir/go-jsonapi-server/model"
)

// Related is a handler for requests corresponding to a resource's
// related resources (primary resource of type t string and identifier
// i string with relationship keyed by k string), with possible
// methods GET and HEAD.
func (env *Environment) HandleRelated(w http.ResponseWriter, r *http.Request) {

	var response Response = NewResponse()

	q, e := model.ParseQueryString(r.URL, env.Parameters)
	if e != nil {
		env.Fail(w, r, e)
		return
	}

	t := chi.URLParam(r, "type")
	i := chi.URLParam(r, "id")
	k := chi.URLParam(r, "related")

	switch r.Method {

	case "OPTIONS":

		_, e := model.GetRelated(r.Context(), env.Graph, t, i, k, env.BaseURL, q)
		if e != nil {
			env.Fail(w, r, e)
			return
		}
		response.Header.Set("Allow", "OPTIONS, GET, HEAD")
		response.Header.Set("Access-Control-Allow-Methods", "OPTIONS, GET, HEAD")
		response.Status = http.StatusNoContent
		env.Success(w, r, response)
		return

	case "GET", "HEAD":

		document, e := model.GetRelated(r.Context(), env.Graph, t, i, k, env.BaseURL, q)
		if e != nil {
			env.Fail(w, r, e)
			return
		}
		response.Body = document
		response.Status = http.StatusOK
		env.Success(w, r, response)
		return

	default:

		// HTTP Method not allowed
		e := model.MakeError(http.StatusMethodNotAllowed)
		e.Code = "ce3a82"
		env.Fail(w, r, e)
		return

	}
}
