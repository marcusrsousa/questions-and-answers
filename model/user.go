package model

import (
	"time"
)

type User struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Login     string     `json:"login,omitempty"`
	Password  string     `json:"password,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}
