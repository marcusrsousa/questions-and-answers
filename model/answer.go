package model

import (
	"time"
)

type Answer struct {
	ID         uint       `gorm:"primary_key" json:"id,omitempty"`
	Text       string     `json:"text,omitempty"`
	Likes      uint       `json:"likes,omitempty"`
	QuestionID uint       `json:"question_id,omitempty"`
	Question   *Question  `json:"question,omitempty"`
	UserID     uint       `json:"user_id,omitempty"`
	User       *User      `json:"user,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `sql:"index" json:"deleted_at,omitempty"`
}
