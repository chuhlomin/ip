package main

import "net/http"

func RealIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check `X-Forwarded-For` header
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			r.RemoteAddr = xff
		}

		next.ServeHTTP(w, r)
	})
}
