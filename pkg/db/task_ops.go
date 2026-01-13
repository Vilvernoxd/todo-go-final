package db

import (
	"database/sql"
	"fmt"
)

func GetTask(id string) (*Task, error) {
	var t Task
	err := db.QueryRow(
		`SELECT CAST(id AS TEXT), date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func UpdateTask(task *Task) error {
	res, err := db.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
