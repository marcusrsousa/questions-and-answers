package model

import (
	"time"
)

type Question struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty"`
	Statement string     `json:"statement,omitempty"`
	Answers   []Answer   `json:"answers,omitempty"`
	UserID    uint       `json:"user_id,omitempty"`
	User      *User      `json:"user,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}
