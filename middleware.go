package main

import "net/http"

func logging(env *environment) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				env.Stdout.Println(
					r.Method,
					r.Proto,
					r.URL.Path,
					r.RemoteAddr,
					r.UserAgent(),
				)
			}()

			next.ServeHTTP(w, r)

		})
	}
}
