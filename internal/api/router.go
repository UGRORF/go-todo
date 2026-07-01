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
	r.Post("/api/signin", SinginHandler)

	r.Group(func(r chi.Router) {
		r.Use(auth)

		r.Post("/api/task", CreateTask)
		r.Get("/api/tasks", GetTasks)
		r.Get("/api/task", GetTask)
		r.Put("/api/task", UpdateTask)
		r.Post("/api/task/done", DoneTask)
		r.Delete("/api/task", DeleteTask)
	})

	return r
}
