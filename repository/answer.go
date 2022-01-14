package repository

import (
	"level7/questions-and-answers/model"

	"github.com/jinzhu/gorm"
)

type AnswerRepository struct {
	Db *gorm.DB
}

func (ar *AnswerRepository) Insert(a *model.Answer) {
	ar.Db.Create(a)
}

func (ar *AnswerRepository) Update(answer *model.Answer, a *model.Answer) {
	ar.Db.Model(answer).Select("Text").Updates(a)
}

func (ar *AnswerRepository) Delete(a *model.Answer) {
	ar.Db.Delete(a)
}

func (ar *AnswerRepository) FindByQuestionId(id uint) *model.Question {
	question := &model.Question{}
	ar.Db.Preload("Answers").First(question, id)
	return question
}

func (ar *AnswerRepository) FindById(id uint) *model.Answer {
	answer := &model.Answer{}
	ar.Db.Select(publicFields).First(answer, id)
	return answer
}
