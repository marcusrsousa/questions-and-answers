package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"level7/questions-and-answers/model"
	"level7/questions-and-answers/repository"

	"github.com/gorilla/mux"
)

type QuestionController struct {
	Repository *repository.QuestionRepository
}

func (qc *QuestionController) Add(w http.ResponseWriter, req *http.Request, currentUser model.User) {
	question := getQuestionFromBody(req)
	defer req.Body.Close()

	question.UserID = currentUser.ID
	question.User = &currentUser

	qc.Repository.Insert(question)

	question.User = nil

	writeResponse(&w, http.StatusCreated, question)
}

func (qc *QuestionController) Update(w http.ResponseWriter, req *http.Request, currentUser model.User) {

	id, err := getUintField(req, "id")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	question := getQuestionFromBody(req)
	defer req.Body.Close()

	if question.ID == 0 {
		question.ID = id
	}

	if question.ID != id {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalQuestion := qc.Repository.FindById(id)

	if originalQuestion.ID == 0 || (originalQuestion.UserID != currentUser.ID && originalQuestion.Statement != question.Statement) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qc.Repository.Update(originalQuestion, question)

	question.User = nil

	writeResponse(&w, http.StatusAccepted, question)

}

func (qc *QuestionController) Delete(w http.ResponseWriter, req *http.Request, currentUser model.User) {

	id, err := getUintField(req, "id")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	question := qc.Repository.FindById(id)

	if question.ID != id || currentUser.ID != question.UserID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qc.Repository.Delete(question)

	writeResponse(&w, http.StatusNoContent, nil)

}

func (qc *QuestionController) GetById(w http.ResponseWriter, req *http.Request, currentUser model.User) {

	id, errId := getUintField(req, "id")

	if errId != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	question := qc.Repository.FindById(id)

	if question.ID == 0 || question.UserID != currentUser.ID {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	writeResponse(&w, http.StatusOK, question)

}

func (qc *QuestionController) Get(w http.ResponseWriter, req *http.Request) {

	user, errUser := getUintField(req, "user_id")

	if errUser != nil {
		writeResponse(&w, http.StatusOK, qc.Repository.FindAll())
		return
	}

	writeResponse(&w, http.StatusOK, qc.Repository.FindByUser(user))

}

func getQuestionFromBody(req *http.Request) *model.Question {
	question := &model.Question{}
	b, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Fatalln("Error reading body:", err)
	}
	json.Unmarshal(b, question)

	return question
}

func getUintField(req *http.Request, field string) (uint, error) {
	id64, err := strconv.ParseUint(mux.Vars(req)[field], 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id64), nil
}

func writeResponse(resp *http.ResponseWriter, statusCode int, v interface{}) {
	w := *resp
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)

	resp = &w
}
