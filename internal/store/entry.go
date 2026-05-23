package store

import "time"

type Entry struct {
	ID        int64
	Title     string
	Content   string
	Category  string
	Tags      string
	FilePath  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
