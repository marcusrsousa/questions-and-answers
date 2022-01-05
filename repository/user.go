package repository

import (
	"level7/questions-and-answers/model"

	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	Db *gorm.DB
}

func (ur *UserRepository) Insert(u *model.User) {
	ur.Db.Create(u)
}

func (ur *UserRepository) FindById(id uint) *model.User {
	return findUserById(id, ur.Db)
}

func (ur *UserRepository) Login(login string) *model.User {
	user := &model.User{}
	ur.Db.Select([]string{"id", "name", "password"}).First(user, "login = ?", login)
	return user
}

func findUserById(id uint, db *gorm.DB) *model.User {
	user := &model.User{}
	db.First(user, id)
	return user
}
