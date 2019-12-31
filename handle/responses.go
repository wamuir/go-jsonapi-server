package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/wamuir/go-jsonapi-server/model"
	"github.com/wamuir/go-jsonapi-server/schema"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"
)

// Response is the header, body and status for a response.
type Response struct {
	Header http.Header
	Body   *model.Document
	Status int
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
func validateMIME(contentType string) *model.ErrorObject {

	//contentType := r.Header.Get("Content-Type")
	mediatype, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		errObj := model.MakeError(http.StatusInternalServerError)
		errObj.Code = "811c93"
		errObj.Title = "Encountered internal error while parsing media type"
		errObj.Detail = err.Error()
		return errObj
	} else if mediatype != "application/vnd.api+json" || len(params) != 0 {
		errObj := model.MakeError(http.StatusUnsupportedMediaType)
		errObj.Code = "d24289"
		errObj.Title = "Invalid media type"
		errObj.Detail = "Clients MUST send all JSON:API data in request documents with the header `Content-Type: application/vnd.api+json` without any media type parameters"
		return errObj
	}

	return nil
}

// Fail completes an errored response.
func Fail(ctx context.Context, stderr *log.Logger, w http.ResponseWriter, r *http.Request, start time.Time, errObj *model.ErrorObject) {

	//var document model.Document
	var document model.Document = model.NewDocument()
	document.Errors = append(document.Errors, *errObj)

	status, err := strconv.Atoi(errObj.Status)
	if err != nil {
		status = http.StatusInternalServerError
		errObj = model.MakeError(http.StatusInternalServerError)
		errObj.Code = "dcaa49"
		errObj.Title = "Encountered an internal error while processing an error: failed to convert HTTP status string to an integer"
		errObj.Detail = err.Error()
		document.Errors = append(document.Errors, *errObj)
	}

	if ctx.Err() != nil {
		status = http.StatusGatewayTimeout
		errObj = model.MakeError(http.StatusGatewayTimeout)
		errObj.Code = "95fd64"
		errObj.Title = "Server timed out while completing the request"
		errObj.Detail = ctx.Err().Error()
		stderr.Println(errObj.Detail)
		document.Errors = append([]model.ErrorObject{*errObj}, document.Errors...) // prepend
	}

	result, err := schema.Validate(document)
	if err != nil {
		errObj := model.MakeError(http.StatusInternalServerError)
		errObj.Code = "f42960"
		errObj.Title = "Encountered an internal error while processing an error: JSON schema validator returned error"
		errObj.Detail = err.Error()
		document.Errors = append(document.Errors, *errObj)
	} else if !result.Valid() {
		errObj := model.MakeError(http.StatusInternalServerError)
		errObj.Code = "b6be49"
		errObj.Title = "Encountered an internal error while processing an error: error document failed to validate against JSON:API schema"
		// errObj.Detail =
		document.Errors = append(document.Errors, *errObj)
	}

	if status == http.StatusInternalServerError {
		defer func() {
			dump := spew.Sdump(errObj)
			stderr.Println(dump)
		}()
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	duration := fmt.Sprintf("%.4f", time.Since(start).Seconds())
	w.Header().Set("Server-Timing", "total;dur="+duration)

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
func Success(ctx context.Context, stderr *log.Logger, w http.ResponseWriter, r *http.Request, start time.Time, response Response) {

	// Validate response body against JSON:API schema
	if response.Body != nil {

		// Set JSON:API Version
		response.Body.Version()

		// Validate
		result, err := schema.Validate(response.Body)
		if err != nil {
			errObj := model.MakeError(http.StatusInternalServerError)
			errObj.Code = "612382"
			errObj.Title = "Encountered internal error while validating request body against JSON:API schema"
			errObj.Detail = err.Error()
			Fail(ctx, stderr, w, r, start, errObj)
			return
		} else if !result.Valid() {
			errObj := model.MakeError(http.StatusInternalServerError)
			errObj.Code = "3d2798"
			errObj.Title = "Response document failed to validate against JSON:API schema"
			// errObj.Detail =
			Fail(ctx, stderr, w, r, start, errObj)
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

	// Add timing to header
	duration := fmt.Sprintf("%.4f", time.Since(start).Seconds())
	response.Header.Set("Server-Timing", "total;dur="+duration)

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
