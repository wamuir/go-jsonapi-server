package handle

import (
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/go-chi/chi"
	"github.com/wamuir/go-jsonapi-core"
	"github.com/wamuir/go-jsonapi-server/model"
)

// Collection is a handler for requests corresponding to a collection of
// resources (of type t string), with possible methods GET, HEAD and POST.
func (env *Environment) HandleCollection(w http.ResponseWriter, r *http.Request) {

	var (
		response Response  = NewResponse()
		start    time.Time = time.Now()
	)

	t := chi.URLParam(r, "type")

	q, e := model.ParseQueryString(r.URL, env.Parameters)
	if e != nil {
		env.Fail(w, r, e)
		return
	}

	switch r.Method {

	case "OPTIONS":

		// Verify collection exists
		_, e := model.GetCollection(r.Context(), env.Graph, t, env.BaseURL, q)
		if e != nil {
			env.Fail(w, r, e)
			return
		}
		response.Header.Set("Allow", "OPTIONS, GET, HEAD, POST")
		response.Header.Set("Access-Control-Allow-Methods", "OPTIONS, GET, HEAD, POST")
		response.Status = http.StatusNoContent
		env.Success(w, r, response)
		return

	case "GET", "HEAD":

		// Get collection
		document, e := model.GetCollection(r.Context(), env.Graph, t, env.BaseURL, q)
		if e != nil {
			env.Fail(w, r, e)
			return
		}
		document.Meta["took"] = time.Now().Sub(start).Milliseconds()
		response.Body = document
		response.Status = http.StatusOK
		env.Success(w, r, response)
		return

	case "POST":

		// Validate content type
		e := ValidateMIME(r.Header.Get("Content-Type"))
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

		// Post new resource
		i, e := model.PostResource(r.Context(), env.Graph, t, document)
		if e != nil {
			env.Fail(w, r, e)
			return
		}

		// Get the resource
		document, e = model.GetResource(r.Context(), env.Graph, i.Type, i.Identifier, env.BaseURL, q)
		if e != nil {
			env.Fail(w, r, e)
			return
		}

		// Build link for the new resource
		ref, err := url.Parse(
			path.Join(i.Type, i.Identifier),
		)
		if err != nil {
			e := core.MakeError(http.StatusInternalServerError)
			e.Code = "bfe23f"
			e.Title = "Encountered internal error while generating response"
			e.Detail = err.Error()
			env.Fail(w, r, e)
			return
		}
		document.Meta["took"] = time.Now().Sub(start).Milliseconds()
		response.Body = document
		response.Header.Set("Location", env.BaseURL.ResolveReference(ref).String())
		response.Status = http.StatusCreated
		env.Success(w, r, response)
		return

	default:

		// HTTP Method not allowed
		e := core.MakeError(http.StatusMethodNotAllowed)
		e.Code = "8e5fce"
		env.Fail(w, r, e)
		return

	}
}
