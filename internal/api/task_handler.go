package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/UGRORF/go-todo/pkg/db"
)

const taskLimit = 50

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &task)
	if err != nil {
		writeError(w, "invalid json body format", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "Title is not defined", http.StatusBadRequest)
		return
	}

	err = CheckDate(&task)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{"id": id}, http.StatusCreated)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []*db.Task
	var err error
	search := r.URL.Query().Get("search")
	if search == "" {
		tasks, err = db.Tasks(taskLimit)
	} else {
		tasks, err = db.SearchTasks(taskLimit, search)
	}

	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, TaskResp{Tasks: tasks}, http.StatusOK)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "id is not defined", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, task, http.StatusOK)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeError(w, "id is required", http.StatusBadRequest)
		return
	}

	_, err = db.GetTask(task.ID)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	if task.Title == "" {
		writeError(w, "title is not defined", http.StatusBadRequest)
		return
	}

	err = CheckDate(&task)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func DoneTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "id is not defined", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(task.ID)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.UpdateDate(task)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "id is not defined", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.DeleteTask(task.ID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if data == nil {
		w.Write([]byte("{}"))
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func writeError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, map[string]string{"error": message}, status)
}
