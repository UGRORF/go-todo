package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RouterInit() *chi.Mux {
	r := chi.NewRouter()

	webDir := "./web"
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(webDir))))
	r.Get("/api/nextdate", GetNextDate)
	return r
}
