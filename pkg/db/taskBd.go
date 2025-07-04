package db

import (
	"errors"
)

func Tasks(limit int) ([]*Task, error) {
	var tasks []*Task
	query := "SELECT * FROM scheduler ORDER BY date LIMIT ?"
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, errors.New("Error Querying DB")
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id      int64
			date    string
			title   string
			comment string
			repeat  string
		)
		err := rows.Scan(&id, &date, &title, &comment, &repeat)
		if err != nil {
			return nil, errors.New("error scanning task")
		}
		task := &Task{
			ID:      id,
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
