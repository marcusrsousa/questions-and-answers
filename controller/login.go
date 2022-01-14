package controller

import (
	"encoding/json"
	"io/ioutil"
	"level7/questions-and-answers/model"
	"level7/questions-and-answers/repository"
	"level7/questions-and-answers/utils"
	"log"
	"net/http"
)

type LoginController struct {
	Repository *repository.UserRepository
}

func (lc *LoginController) Login(w http.ResponseWriter, req *http.Request) {
	userRequest := getUserFromBody(req)

	if userRequest.Login == "" || userRequest.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := lc.Repository.Login(userRequest.Login)

	if user.ID == 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !utils.IsSamePasswords(user.Password, userRequest.Password) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	user.Password = ""

	token, err := utils.GenerateJWT(*user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(token)

}

func (lc *LoginController) SignUp(w http.ResponseWriter, req *http.Request) {
	userRequest := getUserFromBody(req)

	if userRequest.Name == "" || userRequest.Login == "" || userRequest.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userRequest.Password = utils.EncryptPassword(userRequest.Password)

	user := lc.Repository.Login(userRequest.Login)

	if user.ID != 0 {
		w.WriteHeader(http.StatusConflict)
		return
	}

	lc.Repository.Update(userRequest, user)

}

func (lc *LoginController) ChangePassword(w http.ResponseWriter, req *http.Request, currentUser model.User) {
	passwords := &model.Password{}
	b, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Fatalln("Error reading body:", err)
	}
	json.Unmarshal(b, passwords)

	if passwords.CurrentPassword == "" || passwords.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := lc.Repository.Login(currentUser.Login)

	if user.ID == 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !utils.IsSamePasswords(user.Password, passwords.CurrentPassword) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userRequest := *user
	userRequest.Password = utils.EncryptPassword(passwords.Password)

	lc.Repository.Update(user, &userRequest)

}

func getUserFromBody(req *http.Request) *model.User {
	user := &model.User{}
	b, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Fatalln("Error reading body:", err)
	}
	json.Unmarshal(b, user)

	return user
}
