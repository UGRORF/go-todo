package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func GetMainPage(w http.ResponseWriter, r *http.Request) {

	// Получаем путь из контекста chi
	path := chi.URLParam(r, "*")
	if path == "" {
		path = "index.html"
	}

	// Формируем полный путь к файлу
	filePath := filepath.Join("/Users/artem/GolandProjects/go-todo/web/", path)

	// Отдаём файл
	http.ServeFile(w, r, filePath)
}
