package model

import (
	"time"
)

type Page struct {
	ID         string    `json:"id"`
	OwnerName  string    `json:"owner_name"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	IsDeleted  bool      `json:"is_deleted"`
}
