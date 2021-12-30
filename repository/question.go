package repository

import (
	"level7/questions-and-answers/model"

	"github.com/jinzhu/gorm"
)

type QuestionRepository struct {
	Db *gorm.DB
}

var publicFields = []string{"id", "statement", "answer", "user_id"}

func (qr *QuestionRepository) Insert(q *model.Question) {
	qr.Db.Create(q)
}

func (qr *QuestionRepository) Update(question *model.Question, q *model.Question) {
	qr.Db.Model(question).Select("Statement", "Answer", "UserID").Updates(q)
}

func (qr *QuestionRepository) Delete(q *model.Question) {
	qr.Db.Delete(q)
}

func (qr *QuestionRepository) FindAll() *[]model.Question {
	questions := &[]model.Question{}
	qr.Db.Select(publicFields).Find(questions)
	return questions
}

func (qr *QuestionRepository) FindById(id uint) *model.Question {
	return findById(id, qr.Db)
}

func (qr *QuestionRepository) FindByUser(userId uint) *[]model.Question {
	questions := &[]model.Question{}
	qr.Db.Select(publicFields).Find(questions, "user_id = ?", userId)
	return questions
}

func findById(id uint, db *gorm.DB) *model.Question {
	question := &model.Question{}
	db.Select(publicFields).First(question, id)
	return question
}
