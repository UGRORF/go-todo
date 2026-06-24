package api

import (
	"net/http"
	"time"
)

func GetNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		now, _ = time.Parse("20060102", nowStr)
	}

	nextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if nextDate == "" {
		w.Write([]byte(""))
		return
	}

	w.Write([]byte(nextDate))
}
