package main

import (
	"net/http"
	"slices"
)

var HTTPMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodHead,
	http.MethodOptions,
	http.MethodTrace,
	http.MethodConnect,
}

type Decorator func(http.HandlerFunc) http.HandlerFunc

// RequireHTTPMethods ensures that the request has only the specified HTTP Method verb.
// Returns http.StatusMethodNotAllowed if not.
func RequireHTTPMethods(methods ...string) Decorator {
	for _, method := range methods {
		if slices.Index(HTTPMethods, method) < 0 {
			panic("Invalid HTTP method: " + method)
		}
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if slices.Index(methods, r.Method) == -1 {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}
