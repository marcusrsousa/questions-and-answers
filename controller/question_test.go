package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"level7/questions-and-answers/controller"
	"level7/questions-and-answers/model"
	"level7/questions-and-answers/repository"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func TestAdd(t *testing.T) {
	const newId = 1
	question := model.Question{
		Statement: "what is the best ORM in goLang",
		Answer:    "gORM is the best!",
		User:      &model.User{ID: 1},
	}

	qc, mock, err := getQuestionController()
	if err != nil {
		t.Errorf(err.Error())
	}

	defer qc.Repository.Db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "questions" ("statement","answer","user_id","created_at","updated_at","deleted_at") 
			VALUES ($1,$2,$3,$4,$5,$6) RETURNING "questions"."id"`)).
		WithArgs(question.Statement, question.Answer, question.User, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newId))
	mock.ExpectCommit()

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(question)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/question", &buf)
	w := httptest.NewRecorder()

	qc.Add(w, req, *question.User)

	result := w.Result()

	defer result.Body.Close()

	newQuestion := model.Question{}

	json.NewDecoder(result.Body).Decode(&newQuestion)

	if question.Statement != newQuestion.Statement || question.Answer != newQuestion.Answer || question.User.ID != newQuestion.User.ID || newId != newQuestion.ID {
		t.Errorf("expected body: %v, got: %v", question, newQuestion)
	}

	if result.StatusCode != http.StatusCreated {
		t.Errorf("expected status code: %v, got: %v", http.StatusCreated, result.StatusCode)
	}

}

func TestUpdate(t *testing.T) {
	question := model.Question{
		ID:        1,
		Statement: "what is the best ORM in goLang",
		Answer:    "gORM is the best!",
		User:      &model.User{ID: 1},
	}

	qc, mock, err := getQuestionController()
	if err != nil {
		t.Errorf(err.Error())
	}

	defer qc.Repository.Db.Close()

	rows := sqlmock.
		NewRows([]string{"id", "statement", "answer", "user_id", "created_at", "updated_at", "deleted_at"}).
		AddRow(question.ID, question.Statement, question.Answer, question.User, question.CreatedAt, question.UpdatedAt, question.DeletedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, statement, answer, user_id FROM "questions"  WHERE "questions"."deleted_at" IS NULL AND (("questions"."id" = 1)) ORDER BY "questions"."id" ASC LIMIT 1`)).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "questions"`)).WillReturnResult(sqlmock.NewResult(int64(question.ID), 1))
	mock.ExpectCommit()

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(question)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPut, "/question", &buf)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	qc.Update(w, req, *question.User)

	result := w.Result()

	defer result.Body.Close()

	if result.StatusCode != http.StatusAccepted {
		t.Errorf("expected status code: %v, got: %v", http.StatusAccepted, result.StatusCode)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestDelete(t *testing.T) {
	question := model.Question{
		ID:        1,
		Statement: "what is the best ORM in goLang",
		Answer:    "gORM is the best!",
		User:      &model.User{ID: 1},
	}

	qc, mock, err := getQuestionController()
	if err != nil {
		t.Errorf(err.Error())
	}

	defer qc.Repository.Db.Close()

	rows := sqlmock.
		NewRows([]string{"id", "statement", "answer", "user_id", "created_at", "updated_at", "deleted_at"}).
		AddRow(question.ID, question.Statement, question.Answer, question.User, question.CreatedAt, question.UpdatedAt, question.DeletedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, statement, answer, user_id FROM "questions"  WHERE "questions"."deleted_at" IS NULL AND (("questions"."id" = 1)) ORDER BY "questions"."id" ASC LIMIT 1`)).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "questions"`)).WillReturnResult(sqlmock.NewResult(int64(question.ID), 1))
	mock.ExpectCommit()

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(question)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/question", &buf)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	qc.Delete(w, req, *question.User)

	result := w.Result()

	defer result.Body.Close()

	if result.StatusCode != http.StatusNoContent {
		t.Errorf("expected status code: %v, got: %v", http.StatusNoContent, result.StatusCode)
	}

}

func TestGetById(t *testing.T) {
	now := time.Now()
	question := model.Question{
		ID:        1,
		Statement: "what is the best ORM in goLang",
		Answer:    "gORM is the best!",
		User:      &model.User{ID: 1},
		UpdatedAt: &now,
	}

	qc, mock, err := getQuestionController()
	if err != nil {
		t.Errorf(err.Error())
	}

	defer qc.Repository.Db.Close()

	rows := sqlmock.
		NewRows([]string{"id", "statement", "answer", "user_id", "created_at", "updated_at", "deleted_at"}).
		AddRow(question.ID, question.Statement, question.Answer, question.User, question.CreatedAt, question.UpdatedAt, question.DeletedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, statement, answer, user_id FROM "questions"  WHERE "questions"."deleted_at" IS NULL AND (("questions"."id" = 1)) ORDER BY "questions"."id" ASC LIMIT 1`)).
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/question", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}
	req = mux.SetURLVars(req, vars)

	qc.Get(w, req)

	result := w.Result()

	defer result.Body.Close()

	newQuestion := model.Question{}

	json.NewDecoder(result.Body).Decode(&newQuestion)

	if question.Statement != newQuestion.Statement || question.Answer != newQuestion.Answer || question.User.ID != newQuestion.User.ID || question.ID != newQuestion.ID {
		t.Errorf("expected body: %v, got: %v", question, newQuestion)
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status code: %v, got: %v", http.StatusAccepted, result.StatusCode)
	}

}

func TestGetByUser(t *testing.T) {
	now := time.Now()
	question := model.Question{
		ID:        1,
		Statement: "what is the best ORM in goLang",
		Answer:    "gORM is the best!",
		User:      &model.User{ID: 1},
		UpdatedAt: &now,
	}

	qc, mock, err := getQuestionController()
	if err != nil {
		t.Errorf(err.Error())
	}

	defer qc.Repository.Db.Close()

	rows := sqlmock.
		NewRows([]string{"id", "statement", "answer", "user_id", "created_at", "updated_at", "deleted_at"}).
		AddRow(question.ID, question.Statement, question.Answer, question.User, question.CreatedAt, question.UpdatedAt, question.DeletedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, statement, answer, user_id FROM "questions"  WHERE "questions"."deleted_at" IS NULL AND ((user_id = $1))`)).
		WithArgs(question.User.ID).
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/question", nil)
	w := httptest.NewRecorder()

	values := req.URL.Query()
	values.Add("user", strconv.FormatUint(uint64(question.User.ID), 10))
	req.URL.RawQuery = values.Encode()

	qc.Get(w, req)

	result := w.Result()

	defer result.Body.Close()

	questions := []model.Question{}

	json.NewDecoder(result.Body).Decode(&questions)

	if question.Statement != questions[0].Statement || question.Answer != questions[0].Answer || question.User.ID != questions[0].User.ID || question.ID != questions[0].ID {
		t.Errorf("expected body: %v, got: %v", question, questions)
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status code: %v, got: %v", http.StatusAccepted, result.StatusCode)
	}

}

func TestGetAll(t *testing.T) {
	now := time.Now()
	question := model.Question{
		ID:        1,
		Statement: "what is the best ORM in goLang",
		Answer:    "gORM is the best!",
		User:      &model.User{ID: 1},
		UpdatedAt: &now,
	}

	qc, mock, err := getQuestionController()
	if err != nil {
		t.Errorf(err.Error())
	}

	defer qc.Repository.Db.Close()

	rows := sqlmock.
		NewRows([]string{"id", "statement", "answer", "user_id", "created_at", "updated_at", "deleted_at"}).
		AddRow(question.ID, question.Statement, question.Answer, question.User.ID, question.CreatedAt, question.UpdatedAt, question.DeletedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, statement, answer, user_id FROM "questions"  WHERE "questions"."deleted_at" IS NULL`)).
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/question", nil)
	w := httptest.NewRecorder()

	qc.Get(w, req)

	result := w.Result()

	defer result.Body.Close()

	questions := []model.Question{}

	json.NewDecoder(result.Body).Decode(&questions)

	if question.Statement != questions[0].Statement || question.Answer != questions[0].Answer || question.User.ID != questions[0].User.ID || question.ID != questions[0].ID {
		t.Errorf("expected body: %v, got: %v", question, questions)
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status code: %v, got: %v", http.StatusAccepted, result.StatusCode)
	}

}

func getQuestionController() (*controller.QuestionController, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()

	if err != nil {
		return nil, nil, fmt.Errorf("error while open mock database: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)

	if err != nil {
		return nil, nil, fmt.Errorf("error while open gorm database: %v", err)
	}

	rep := &repository.QuestionRepository{Db: gdb}
	return &controller.QuestionController{Repository: rep}, mock, nil
}
