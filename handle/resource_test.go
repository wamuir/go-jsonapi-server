package handle

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/wamuir/go-jsonapi-server/config"
	sqlite3 "github.com/wamuir/go-jsonapi-server/graph/sqlite3"
)

func TestHandleResource(t *testing.T) {

	var (
		b   []byte
		ctx *chi.Context
		r   *http.Request
		o   *http.Response
		w   *httptest.ResponseRecorder
	)

	/////////////////////////////////////// SETUP

	// graph
	g, err := sqlite3.Connect("file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer g.Close()

	// open /dev/null for logging to nowhere
	devnull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer devnull.Close()

	// set up environment
	e := &Environment{
		Graph:      g,
		Parameters: config.Parameters,
		Stderr:     log.New(devnull, "", 0),
		Stdout:     log.New(devnull, "", 0),
	}

	// define a resource to be posted
	b = []byte(`{"data":{"type":"foo","id":"bar","attributes":{"a":"b"},"meta":{"c":"d"}}}`)

	// create http request
	r = httptest.NewRequest(http.MethodPost, "/foo/", bytes.NewBuffer(b))
	r.Header.Set("Content-Type", "application/vnd.api+json")
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	// post resource
	w = httptest.NewRecorder()
	e.HandleCollection(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(o.Body)
		log.Println(string(b))
		t.Fatalf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusCreated,
		)
	}

	/////////////////////////////////////// TESTS

	// invalid query string
	r = httptest.NewRequest(http.MethodGet, "/?foo=bar,baz", nil)
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusBadRequest,
		)
	}

	// OPTIONS non-existent resource
	r = httptest.NewRequest(http.MethodOptions, "/foo/baz", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	ctx.URLParams.Add("id", "baz")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusNotFound {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusNotFound,
		)
	}

	// OPTIONS resource
	r = httptest.NewRequest(http.MethodOptions, "/foo/bar", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	ctx.URLParams.Add("id", "bar")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusNoContent {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusNoContent,
		)
	}

	// GET non-existent resource
	r = httptest.NewRequest(http.MethodGet, "/foo/baz", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	ctx.URLParams.Add("id", "baz")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusNotFound {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusNotFound,
		)
	}

	// GET resource
	r = httptest.NewRequest(http.MethodGet, "/foo/bar", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	ctx.URLParams.Add("id", "bar")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusOK {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusOK,
		)
	}

	// HEAD non-existent resource
	r = httptest.NewRequest(http.MethodHead, "/foo/baz", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	ctx.URLParams.Add("id", "baz")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusNotFound {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusNotFound,
		)
	}

	// HEAD resource
	r = httptest.NewRequest(http.MethodHead, "/foo/bar", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	ctx.URLParams.Add("id", "bar")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusOK {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusOK,
		)
	}

	// PATCH non-existent resource
	r = httptest.NewRequest(http.MethodPatch, "/foo/baz", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	ctx.URLParams.Add("id", "baz")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusMethodNotAllowed,
		)
	}

	// PATCH resource
	r = httptest.NewRequest(http.MethodPatch, "/foo/bar", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	ctx.URLParams.Add("id", "bar")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusMethodNotAllowed,
		)
	}

	// DELETE non-existent resource
	r = httptest.NewRequest(http.MethodHead, "/foo/baz", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	ctx.URLParams.Add("id", "baz")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusNotFound {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusNotFound,
		)
	}

	// DELETE resource
	r = httptest.NewRequest(http.MethodHead, "/foo/bar", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	ctx.URLParams.Add("id", "bar")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusOK {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusOK,
		)
	}

	// Unsupported method
	r = httptest.NewRequest(http.MethodConnect, "/foo/bar", nil)
	r.Header.Set("Accept", "application/vnd.api+json")
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("type", "foo")
	ctx.URLParams.Add("id", "bar")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	w = httptest.NewRecorder()
	e.HandleResource(w, r)
	o = w.Result()
	if o.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf(
			"o.StatusCode = %v, want %v",
			o.StatusCode,
			http.StatusMethodNotAllowed,
		)
	}

}
