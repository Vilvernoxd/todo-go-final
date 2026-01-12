package db

func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query(`SELECT CAST(id AS TEXT), date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		t := new(Task)
		err = rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
