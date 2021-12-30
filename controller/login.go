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
	user := lc.Repository.Login(userRequest.Login, userRequest.Password)

	if user.ID == 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	token, err := utils.GenerateJWT(*user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(token)

}

func (lc *LoginController) SignUp(w http.ResponseWriter, req *http.Request) {
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
