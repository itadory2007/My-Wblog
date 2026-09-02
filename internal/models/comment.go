package models

import "time"

type Comment struct {
	ID int64
	PostID int64
	AuthorID int64
	Content string
	CreatedAt time.Time
}

