package question_controller

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

func (qc *QuestionController) Add(w http.ResponseWriter, req *http.Request) {
	question := getQuestionFromBody(req)
	defer req.Body.Close()

	qc.Repository.Insert(question)

	writeResponse(&w, http.StatusCreated, question)
}

func (qc *QuestionController) Update(w http.ResponseWriter, req *http.Request) {

	id, err := getId(req)

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

	if originalQuestion.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qc.Repository.Update(originalQuestion, question)

	writeResponse(&w, http.StatusAccepted, question)

}

func (qc *QuestionController) Delete(w http.ResponseWriter, req *http.Request) {

	id, err := getId(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	question := qc.Repository.FindById(id)

	if question.ID != id {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qc.Repository.Delete(question)

	writeResponse(&w, http.StatusNoContent, nil)

}

func (qc *QuestionController) Get(w http.ResponseWriter, req *http.Request) {

	id, err := getId(req)

	user := req.URL.Query().Get("user")

	if err != nil && user == "" {
		writeResponse(&w, http.StatusOK, qc.Repository.FindAll())
		return
	}
	if user != "" {
		writeResponse(&w, http.StatusOK, qc.Repository.FindByUser(user))
		return
	}

	question := qc.Repository.FindById(id)
	if question.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	writeResponse(&w, http.StatusOK, question)

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

func getId(req *http.Request) (uint, error) {
	id64, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 64)
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
