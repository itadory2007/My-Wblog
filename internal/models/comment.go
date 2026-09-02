package models

import "time"

type Post struct {
	ID int64
	Title string
	Content string
	Image *string
	AuthorID int64
	IsPrivate bool
	CreatedAt time.Time
}
