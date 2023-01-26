package celeritas

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (c *Celeritas) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	if c.Debug {
		mux.Use(middleware.Logger)
	}

	// if my app stop working form panic
	mux.Use(middleware.Recoverer)

	// route
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "testing")
	})
	return mux
}
