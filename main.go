package main

import (
	"fmt"
	"github.com/wamuir/go-jsonapi-server/graph"
	sqlite3 "github.com/wamuir/go-jsonapi-server/graph/sqlite3"
	"github.com/wamuir/go-jsonapi-server/model"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type environment struct {
	BaseURL    url.URL
	Graph      graph.Graph
	Parameters model.Parameters
	Stderr     *log.Logger
	Stdout     *log.Logger
}

func main() {

	stderr := log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.LUTC)
	stdout := log.New(os.Stdout, "INFO: ", log.LstdFlags|log.LUTC)

	g, err := sqlite3.Connect(dsn.String())
	if err != nil {
		stderr.Fatal(err.Error())
	}

	env := &environment{
		BaseURL:    baseURL,
		Graph:      g,
		Parameters: parameters,
		Stderr:     stderr,
		Stdout:     stdout,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", listenAddr, listenPort),
		Handler:      logging(env)(route(ctxTimeout, env)),
		ErrorLog:     stderr,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}

	stderr.Fatal(server.ListenAndServe())
}
