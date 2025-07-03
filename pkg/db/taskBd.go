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
			id      int
			title   string
			comment string
			date    string
			repeat  string
		)
		err := rows.Scan(&id, &title, &comment, &date, &repeat)
		if err != nil {
			return nil, errors.New("error scanning task")
		}
		task := &Task{
			id,
			date,
			title,
			comment,
			repeat,
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
