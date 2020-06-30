package handle

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/wamuir/go-jsonapi-core"
	"github.com/wamuir/go-jsonapi-server/graph"
	"github.com/wamuir/go-jsonapi-server/model"
	"github.com/wamuir/go-jsonapi-server/schema"
)

type Environment struct {
	BaseURL    url.URL
	Graph      graph.Graph
	Parameters model.Parameters
	Stderr     *log.Logger
	Stdout     *log.Logger
}

// Response is the header, body and status for a response.
type Response struct {
	Header  http.Header
	Body    *core.Document
	Status  int
	Trailer http.Header
}

// NewResponse is a Response constructor.
func NewResponse() Response {
	response := Response{
		Header: make(http.Header),
	}
	response.Header.Set("Access-Control-Allow-Origin", "*")
	return response
}

// From net/http/httputil
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// Validates Content-Type header per JSON:API spec
func ValidateMIME(contentType string) *core.Error {

	//contentType := r.Header.Get("Content-Type")

	if len(strings.TrimSpace(contentType)) == 0 {
		e := core.MakeError(http.StatusUnsupportedMediaType)
		e.Code = "d24289"
		e.Title = "Missing media type"
		e.Detail = "Clients MUST send all JSON:API data in request documents with the header `Content-Type: application/vnd.api+json` without any media type parameters"
		return e
	}

	mediatype, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		e := core.MakeError(http.StatusInternalServerError)
		e.Code = "811c93"
		e.Title = "Encountered internal error while parsing media type"
		e.Detail = err.Error()
		return e
	} else if mediatype != "application/vnd.api+json" || len(params) != 0 {
		e := core.MakeError(http.StatusUnsupportedMediaType)
		e.Code = "d24289"
		e.Title = "Invalid media type"
		e.Detail = "Clients MUST send all JSON:API data in request documents with the header `Content-Type: application/vnd.api+json` without any media type parameters"
		return e
	}

	return nil
}

// Fail completes an errored response.
func (env *Environment) Fail(w http.ResponseWriter, r *http.Request, e *core.Error) {

	//var document model.Document
	var document core.Document = core.New()
	document.Errors = append(document.Errors, *e)

	status, err := strconv.Atoi(e.Status)
	if err != nil {
		status = http.StatusInternalServerError
		e := core.MakeError(http.StatusInternalServerError)
		e.Code = "dcaa49"
		e.Title = "Encountered an internal error while processing an error: failed to convert HTTP status string to an integer"
		e.Detail = err.Error()
		document.Errors = append(document.Errors, *e)
	}

	if r.Context().Err() != nil {
		status = http.StatusGatewayTimeout
		e := core.MakeError(http.StatusGatewayTimeout)
		e.Code = "95fd64"
		e.Title = "Server timed out while completing the request"
		e.Detail = r.Context().Err().Error()
		env.Stderr.Println(e.Detail)
		document.Errors = append([]core.Error{*e}, document.Errors...) // prepend
	}

	result, err := schema.Validate(document)
	if err != nil {
		e := core.MakeError(http.StatusInternalServerError)
		e.Code = "f42960"
		e.Title = "Encountered an internal error while processing an error: JSON schema validator returned error"
		e.Detail = err.Error()
		document.Errors = append(document.Errors, *e)
	} else if !result.Valid() {
		e := core.MakeError(http.StatusInternalServerError)
		e.Code = "b6be49"
		e.Title = "Encountered an internal error while processing an error: error document failed to validate against JSON:API schema"
		// errObj.Detail =
		document.Errors = append(document.Errors, *e)
	}

	if status == http.StatusInternalServerError {
		defer func() {
			dump := spew.Sdump(document.Errors)
			env.Stderr.Println(dump)
		}()
	}

	if r.Method == "HEAD" {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	encoder.Encode(&document)

	return
}

// Success completes an unerrored response.
func (env *Environment) Success(w http.ResponseWriter, r *http.Request, response Response) {

	// Validate response body against JSON:API schema
	if response.Body != nil {

		// Set JSON:API Version
		response.Body.Version()

		// Validate
		result, err := schema.Validate(response.Body)
		if err != nil {
			errObj := core.MakeError(http.StatusInternalServerError)
			errObj.Code = "612382"
			errObj.Title = "Encountered internal error while validating request body against JSON:API schema"
			errObj.Detail = err.Error()
			env.Fail(w, r, errObj)
			return

		} else if !result.Valid() {
			errObj := core.MakeError(http.StatusInternalServerError)
			errObj.Code = "3d2798"
			errObj.Title = "Response document failed to validate against JSON:API schema"
			// errObj.Detail =
			env.Fail(w, r, errObj)
			return
		}

		// Set Headers as appropriate given HTTP Method
		if r.Method == "HEAD" {
			length := response.Body.ContentLength()
			response.Header.Set("Content-Length", strconv.Itoa(length))
		} else {
			response.Header.Set("X-Content-Type-Options", "nosniff")
			response.Header.Set("Content-Type", "application/vnd.api+json")
		}
	}

	// Write header
	copyHeader(w.Header(), response.Header)
	w.WriteHeader(response.Status)

	// Write body
	if response.Body != nil && r.Method != "HEAD" {
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "\t")
		encoder.Encode(response.Body)
	}
}
