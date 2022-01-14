package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"level7/questions-and-answers/model"
	"level7/questions-and-answers/repository"

	"github.com/gorilla/mux"
)

type AnswerController struct {
	Repository *repository.AnswerRepository
}

func (ac *AnswerController) Add(w http.ResponseWriter, req *http.Request, currentUser model.User) {
	answer := getAnswerFromBody(req)
	defer req.Body.Close()

	answer.UserID = currentUser.ID
	answer.User = &currentUser

	ac.Repository.Insert(answer)

	answer.User = nil

	writeResponse(&w, http.StatusCreated, answer)
}

func (ac *AnswerController) Update(w http.ResponseWriter, req *http.Request, currentUser model.User) {

	id, err := getUintField(mux.Vars(req)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	answer := getAnswerFromBody(req)
	defer req.Body.Close()

	if answer.ID == 0 {
		answer.ID = id
	}

	if answer.ID != id {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalAnswer := ac.Repository.FindById(id)

	if originalAnswer.ID == 0 || (originalAnswer.UserID != currentUser.ID && originalAnswer.Text != originalAnswer.Text) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ac.Repository.Update(originalAnswer, answer)

	answer.User = nil

	writeResponse(&w, http.StatusAccepted, answer)

}

func (ac *AnswerController) Delete(w http.ResponseWriter, req *http.Request, currentUser model.User) {

	id, err := getUintField(mux.Vars(req)["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	answer := ac.Repository.FindById(id)

	if answer.ID != id || currentUser.ID != answer.UserID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ac.Repository.Delete(answer)

	writeResponse(&w, http.StatusNoContent, nil)

}

func (ac *AnswerController) GetById(w http.ResponseWriter, req *http.Request, currentUser model.User) {

	id, errId := getUintField(mux.Vars(req)["id"])

	if errId != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errId)
		return
	}

	answer := ac.Repository.FindById(id)

	if answer.ID == 0 || answer.UserID != currentUser.ID {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Answer not found to this user")
		return
	}
	writeResponse(&w, http.StatusOK, answer)

}

func (qc *AnswerController) Get(w http.ResponseWriter, req *http.Request) {

	questionId, errId := getUintField(mux.Vars(req)["questionId"])

	if errId != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error to get question id:", errId)
		return
	}

	writeResponse(&w, http.StatusOK, qc.Repository.FindByQuestionId(questionId))

}

func getAnswerFromBody(req *http.Request) *model.Answer {
	questionId, errId := getUintField(mux.Vars(req)["questionId"])

	if errId != nil {
		log.Fatalln("Error to get question id:", errId)
		return nil
	}

	answer := &model.Answer{}
	b, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Fatalln("Error reading body:", err)
		return nil
	}
	json.Unmarshal(b, answer)

	answer.QuestionID = questionId

	return answer
}
