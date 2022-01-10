package repository

import (
	"level7/questions-and-answers/model"
	"level7/questions-and-answers/utils/pagination"

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

func (qr *QuestionRepository) FindAll(page *pagination.Page) pagination.Pagination {
	questions := &[]model.Question{}
	qr.Db.Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select([]string{"id", "name"})
	}).Limit(page.GetLimit()).Offset(page.GetOffset()).Find(questions)

	var count uint64
	qr.Db.Model(&model.Question{}).Count(&count)
	return *pagination.CreatePagination(*page, questions, count)
}

func (qr *QuestionRepository) FindById(id uint) *model.Question {
	return findById(id, qr.Db)
}

func (qr *QuestionRepository) FindByUser(userId uint, page *pagination.Page) pagination.Pagination {
	questions := &[]model.Question{}
	qr.Db.Select(publicFields).Limit(page.GetLimit()).Offset(page.GetOffset()).Find(questions, "user_id = ?", userId)
	var count uint64
	qr.Db.Model(&model.Question{}).Where("user_id = ?", userId).Count(&count)
	return *pagination.CreatePagination(*page, questions, count)
}

func findById(id uint, db *gorm.DB) *model.Question {
	question := &model.Question{}
	db.Select(publicFields).First(question, id)
	return question
}
