package model

import (
	"time"
)

type Page struct {
	ID         string
	OwnerName  string
	Title      string
	Content    string
	CreatedAt  time.Time
	ModifiedAt time.Time
	IsDeleted  bool
}
