package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask вставляет новую задачу и возвращает её ID.
func AddTask(task *Task) (int64, error) {
	const query = `
        INSERT INTO scheduler (date, title, comment, repeat)
        VALUES (?, ?, ?, ?)
    `
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}
	return res.LastInsertId()
}

// Tasks возвращает до limit ближайших задач, отсортированных по дате.
func Tasks(limit int) ([]*Task, error) {
	const query = `
        SELECT id, date, title, comment, repeat
          FROM scheduler
      ORDER BY date
         LIMIT ?
    `
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	list := make([]*Task, 0)
	for rows.Next() {
		var (
			id      int64
			date    string
			title   string
			comment string
			repeat  string
		)
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		list = append(list, &Task{
			ID:      strconv.FormatInt(id, 10),
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return list, nil
}

// GetTask возвращает задачу по её ID.
func GetTask(id string) (*Task, error) {
	const query = `
        SELECT id, date, title, comment, repeat
          FROM scheduler
         WHERE id = ?
    `
	row := DB.QueryRow(query, id)
	var (
		idInt   int64
		date    string
		title   string
		comment string
		repeat  string
	)
	if err := row.Scan(&idInt, &date, &title, &comment, &repeat); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("get task: %w", err)
	}
	return &Task{
		ID:      strconv.FormatInt(idInt, 10),
		Date:    date,
		Title:   title,
		Comment: comment,
		Repeat:  repeat,
	}, nil
}

// UpdateTask обновляет все поля задачи по её ID.
func UpdateTask(task *Task) error {
	const query = `
        UPDATE scheduler
           SET date    = ?,
               title   = ?,
               comment = ?,
               repeat  = ?
         WHERE id = ?
    `
	res, err := DB.Exec(query,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// DeleteTask удаляет задачу по её ID.
func DeleteTask(id string) error {
	const query = `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// UpdateDate меняет только дату задачи по её ID.
func UpdateDate(nextDate, id string) error {
	const query = `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, nextDate, id)
	if err != nil {
		return fmt.Errorf("update date: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
