package config

import (
	"github.com/wamuir/go-jsonapi-server/model"
	"net/url"
	"path/filepath"
)

// The public identity of the server.
//
// This is used to generate JSON:API links and should be the base uri
// of the server, including the scheme (http, https) and the fully-
// qualified domain name.  Optionally, it may include a port and any
// stub for the URI.
//
//  examples:  http://www.example.com:8080/api/
//             https://api.example.com/v1/
//             https://api.example.com/
//
var BaseURL = url.URL{
	Scheme: "http",
	Host:   "localhost:8080", // port optional
	Path:   "/",              // trailing slash required
}

// Address and port for the HTTP server to listen on.
//
//   listenAddr: the address (or hostname) to listen on
//   listenPort: the port number to listen on
//
var (
	ListenAddr string = "127.0.0.1"
	ListenPort int    = 8080
)

// Data source name (DSN) for connection to backend data store.
//
//      sqlite3:  see github.com/mattn/go-sqlite3 for additional info
//        other:  reference the relevant documentation
//
var DSN = url.URL{
	Scheme: "file",
	Opaque: filepath.FromSlash("/tmp/graph.sqlite3"), // <- /var/db
	RawQuery: url.Values{
		"_busy_timeout": []string{"5000"},
		"cache":         []string{"shared"},
		"_foreign_keys": []string{"ON"},
	}.Encode(),
}

// Timeouts for the HTTP server.  All times are in seconds.
//
//         read:  the maximum duration for reading the entire request,
//                including the request body
//        write:  the maximum duration before timing out writes of the
//                response
//         idle:  the maximum amount of time to wait for the next
//                request when keep-alives are enabled
//
var (
	ReadTimeout  int = 5
	WriteTimeout int = 5
	IdleTimeout  int = 5
	CtxTimeout   int = 4
)

// Query parameters.
//
//      include:  for inclusion of resources related to primary data
//                e.g., ?include=contractor,contractor.subcontractors
//  page[limit]:  page size for paginated collections of resources
//                e.g., ?page[limit]=10
// page[offset]:  page offset for paginated collections of resources
//                e.g., ?page[offset]=0
//
var Parameters = model.Parameters{
	"include": model.Parameter{
		Allowed: true,
		Maximum: 3, // Maximum depth for traversal
	},
	"page[limit]": model.Parameter{
		Allowed: true,
		Default: 10,
		Minimum: 1,
		Maximum: 1<<63 - 1,
	},
	"page[offset]": model.Parameter{
		Allowed: true,
		Default: 0,
		Minimum: 0,
		Maximum: 1<<63 - 1,
	},
	"sort": model.Parameter{
		Allowed: true,
	},
}
