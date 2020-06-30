package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/wamuir/go-jsonapi-server/config"
	sqlite3 "github.com/wamuir/go-jsonapi-server/graph/sqlite3"
	"github.com/wamuir/go-jsonapi-server/handle"
)

func main() {

	stderr := log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.LUTC)
	stdout := log.New(os.Stdout, "INFO: ", log.LstdFlags|log.LUTC)

	graph, err := sqlite3.Connect(config.DSN.String())
	if err != nil {
		stderr.Fatal(err.Error())
	}

	env := &handle.Environment{
		BaseURL:    config.BaseURL,
		Graph:      graph,
		Parameters: config.Parameters,
		Stderr:     stderr,
		Stdout:     stdout,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.NoCache)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(config.CtxTimeout) * time.Second))
	r.NotFound(env.Handle404)
	r.MethodNotAllowed(env.Handle405)
	r.Route(`/{type}`, func(r chi.Router) {
		r.HandleFunc("/", env.HandleCollection)
		r.Route(`/{id}`, func(r chi.Router) {
			r.HandleFunc(`/`, env.HandleResource)
			r.HandleFunc(`/{related}`, env.HandleRelated)
			r.HandleFunc(`/relationships/{relationship}`, env.HandleRelationship)
		})
	})

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.ListenAddr, config.ListenPort),
		Handler:      r,
		ErrorLog:     stderr,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.IdleTimeout) * time.Second,
	}

	stderr.Fatal(server.ListenAndServe())
}
