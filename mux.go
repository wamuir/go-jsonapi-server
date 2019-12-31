package main

import (
	"context"
	"github.com/wamuir/go-jsonapi-server/handle"
	"github.com/wamuir/go-jsonapi-server/model"
	"net/http"
	"regexp"
	"time"
)

var (
	isCollection   = regexp.MustCompile(`^\/([^\/]+)\/?$`)
	isResource     = regexp.MustCompile(`^\/([^\/]+)\/([^\/]+)$`)
	isRelationship = regexp.MustCompile(`^\/([^\/]+)\/([^\/]+)\/(?:relationships)\/([^\/]+)\/?$`)
	isRelated      = regexp.MustCompile(`^\/([^\/]+)\/([^\/]+)\/([^\/]+)\/?$`)
)

// Route request to correct handler
func route(ctxTimeout int, env *environment) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Make context
		s := time.Now() // start
		ctx, cancel := context.WithDeadline(
			context.Background(),
			s.Add(time.Duration(ctxTimeout)*time.Second),
		)
		defer cancel()

		// Support clients lacking PATCH
		if r.Method == "POST" && r.Header.Get("X-HTTP-Method-Override") == "PATCH" {
			r.Method = "PATCH"
		}

		switch {

		// Send to collection route
		case isCollection.MatchString(r.URL.Path):
			submatch := isCollection.FindStringSubmatch(r.URL.Path)
			response, errObj := handle.Collection(
				ctx,
				env.BaseURL,
				env.Graph,
				env.Parameters,
				w,
				r,
				submatch[1], // Resource type
			)
			if errObj != nil {
				handle.Fail(ctx, env.Stderr, w, r, s, errObj)
				return
			}
			handle.Success(ctx, env.Stderr, w, r, s, response)

		// Send to resource route
		case isResource.MatchString(r.URL.Path):
			submatch := isResource.FindStringSubmatch(r.URL.Path)
			response, errObj := handle.Resource(
				ctx,
				env.BaseURL,
				env.Graph,
				env.Parameters,
				w,
				r,
				submatch[1], // Resource Type
				submatch[2], // Resource Identifier
			)
			if errObj != nil {
				handle.Fail(ctx, env.Stderr, w, r, s, errObj)
				return
			}
			handle.Success(ctx, env.Stderr, w, r, s, response)

		// Send to relationship route
		case isRelationship.MatchString(r.URL.Path):
			submatch := isRelationship.FindStringSubmatch(r.URL.Path)
			response, errObj := handle.Relationship(
				ctx,
				env.BaseURL,
				env.Graph,
				env.Parameters,
				w,
				r,
				submatch[1], // Resource Type
				submatch[2], // Resource Identifier
				submatch[3], // Relationship Key
			)
			if errObj != nil {
				handle.Fail(ctx, env.Stderr, w, r, s, errObj)
				return
			}
			handle.Success(ctx, env.Stderr, w, r, s, response)

		// Send to related route
		case isRelated.MatchString(r.URL.Path):
			submatch := isRelated.FindStringSubmatch(r.URL.Path)
			response, errObj := handle.Related(
				ctx,
				env.BaseURL,
				env.Graph,
				env.Parameters,
				w,
				r,
				submatch[1], // Resource Type
				submatch[2], // Resource Identifier
				submatch[3], // Relationship Key
			)
			if errObj != nil {
				handle.Fail(ctx, env.Stderr, w, r, s, errObj)
				return
			}
			handle.Success(ctx, env.Stderr, w, r, s, response)

		// Return HTTP 404 (no matching route)
		default:
			errObj := model.MakeError(http.StatusNotFound)
			errObj.Code = "f7519b"
			handle.Fail(ctx, env.Stderr, w, r, s, errObj)
			return

		}

	})
}
