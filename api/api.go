package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type PostBody struct {
    URL string `json:"url"`
}

type Response struct {
    Error string `json:"error,omitempty"`
    Data any `json:"data,omitempty"`
}

func NewHandler() http.Handler {
    router := chi.NewMux()
    // call middlewares
    router.Use(middleware.Recoverer)
    router.Use(middleware.RequestID)
    router.Use(middleware.Logger)

    // set routes
    router.Post("/api/shorten", handleShorten)
    router.Get("/{code}", handleGet)

    return router
}

func handleShorten(rw http.ResponseWriter, req *http.Request) { }

func handleGet(rw http.ResponseWriter, req *http.Request) { }
