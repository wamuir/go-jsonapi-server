package handle

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/wamuir/go-jsonapi-core"
	"github.com/wamuir/go-jsonapi-server/model"
)

// Resource is a handler for requests corresponding to a single resource
// (of type t string with identifier i string), with possible methods
// GET, HEAD, PATCH and DELETE.
func (env *Environment) HandleResource(w http.ResponseWriter, r *http.Request) {

	var (
		response Response  = NewResponse()
		start    time.Time = time.Now()
	)

	q, e := model.ParseQueryString(r.URL, env.Parameters)
	if e != nil {
		env.Fail(w, r, e)
		return
	}

	t := chi.URLParam(r, "type")
	i := chi.URLParam(r, "id")

	switch r.Method {

	case "OPTIONS":

		_, e := model.GetResource(r.Context(), env.Graph, t, i, env.BaseURL, q)
		if e != nil {
			env.Fail(w, r, e)
			return
		}
		response.Header.Set("Allow", "OPTIONS, GET, HEAD, PATCH, DELETE")
		response.Header.Set("Access-Control-Allow-Methods", "OPTIONS, GET, HEAD, PATCH, DELETE")
		response.Status = http.StatusNoContent
		env.Success(w, r, response)
		return

	case "GET", "HEAD":

		document, e := model.GetResource(r.Context(), env.Graph, t, i, env.BaseURL, q)
		if e != nil {
			env.Fail(w, r, e)
			return
		}
		document.Meta["took"] = time.Now().Sub(start).Milliseconds()
		response.Body = document
		response.Status = http.StatusOK
		env.Success(w, r, response)
		return

	case "PATCH":

		// Method not implemented
		e := core.MakeError(http.StatusNotImplemented)
		e.Code = "b83e07"
		env.Success(w, r, response)
		return

	case "DELETE":

		e := model.DeleteResource(r.Context(), env.Graph, t, i)
		if e != nil {
			env.Fail(w, r, e)
			return
		}
		response.Status = http.StatusNoContent
		env.Success(w, r, response)
		return

	default:

		e := core.MakeError(http.StatusMethodNotAllowed)
		e.Code = "594414"
		env.Fail(w, r, e)
		return

	}
}
