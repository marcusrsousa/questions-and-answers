package model

import (
	"time"
)

type Question struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty"`
	Statement string     `json:"statement,omitempty"`
	Answer    string     `json:"answer,omitempty"`
	UserName  string     `json:"user,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}
