package models

import "time"

type Post struct {
	ID          int64    `json:"id"`
	Title        string    `json:"title" validate:"required,min=3,max=255`
	Content        string   `json:"content" validate:"required,min=3,max=500`
	CategoryId int64    `json:"category_id" validate:"required,numeric`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}