package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      string `db:"id" json:"id"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := db1.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}
	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetTask(id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("invalid task id")
	}
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	err := db1.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to query task: %w", err)
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	if task == nil {
		return fmt.Errorf("invalid task")
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err := db1.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil
}

func Tasks(limit int) ([]*Task, error) {
	tasks := make([]*Task, 0)
	query := `SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date ASC, id ASC
			LIMIT ?`

	rows, err := db1.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()
	if rows == nil {
		return []*Task{}, nil
	}

	for rows.Next() {
		task := &Task{}
		rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func SearchTasks(limit int, search string) ([]*Task, error) {
	tasks := make([]*Task, 0)
	query := `SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE ? OR comment LIKE ? OR date LIKE ?
			ORDER BY date ASC, id ASC
			LIMIT ?`

	if isDate(search) {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			return nil, fmt.Errorf("failed to get tasks: %w", err)
		}
		search = date.Format("20060102")
	}

	rows, err := db1.Query(query, "%"+search+"%", "%"+search+"%", "%"+search+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("invalid task id")
	}
	query := `DELETE FROM scheduler WHERE id = ?`
	result, err := db1.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("task with id %s not found", id)
	}
	return nil
}

func UpdateDate(task *Task) error {
	if task == nil {
		return fmt.Errorf("invalid task")
	}
	if task.ID == "" {
		return fmt.Errorf("invalid task id")
	}
	if task.Date == "" {
		return fmt.Errorf("invalid task date")
	}
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	_, err := db1.Exec(query, task.Date, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil
}

func isDate(date string) bool {
	_, err := time.Parse("02.01.2006", date)
	return err == nil
}
