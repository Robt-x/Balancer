package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
)

func MAxAllowedRequests(n uint) mux.MiddlewareFunc {
	queue := make(chan struct{}, n)
	acquire := func() { queue <- struct{}{} }
	release := func() { <-queue }
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acquire()
			defer release()
			next.ServeHTTP(w, r)
		})
	}
}
