package models

import "time"

type Task struct {
	ID          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Status      string `db:"status"`
	UserID      string `db:"userid"`
	UpdateDate  time.Time
}
