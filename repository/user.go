package repository

import (
	"level7/questions-and-answers/model"

	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	Db *gorm.DB
}

var userPublicFields = []string{"id", "name"}

func (ur *UserRepository) Insert(u *model.User) {
	ur.Db.Create(u)
}

func (ur *UserRepository) FindById(id uint) *model.User {
	return findUserById(id, ur.Db)
}

func (ur *UserRepository) Login(login, password string) *model.User {
	user := &model.User{}
	ur.Db.Select(userPublicFields).First(user, "login = ? AND password = ?", login, password)
	return user
}

func findUserById(id uint, db *gorm.DB) *model.User {
	user := &model.User{}
	db.First(user, id)
	return user
}
