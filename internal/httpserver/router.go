package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

type MiddlewareFunc func(http.Handler) http.Handler

//NewRouter returns an HTTP multiplexor
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range declaredRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	router.Use(BasicAuth)
	return router
}
