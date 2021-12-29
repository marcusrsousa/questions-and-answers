package main

import (
	"level7/questions-and-answers/connection"
	question_controller "level7/questions-and-answers/controller"
	"level7/questions-and-answers/model"
	"level7/questions-and-answers/repository"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func addRoutes(r *mux.Router, db *gorm.DB) {

	qc := getQuestionController(db)

	r.HandleFunc("/question", qc.Add).Methods(http.MethodPost)
	r.HandleFunc("/question/{id}", qc.Update).Methods(http.MethodPut)
	r.HandleFunc("/question/{id}", qc.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/question/{id}", qc.Get).Methods(http.MethodGet)
	r.HandleFunc("/question", qc.Get).Methods(http.MethodGet)

}

func getQuestionController(db *gorm.DB) *question_controller.QuestionController {
	rep := &repository.QuestionRepository{Db: db}
	return &question_controller.QuestionController{Repository: rep}
}

func initialMigration(db *gorm.DB) {
	db.AutoMigrate(&model.Question{})
}

func getCors() []handlers.CORSOption {
	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	return []handlers.CORSOption{header, methods, origins}
}

func main() {
	db := connection.GetConnection()
	initialMigration(db)
	r := mux.NewRouter()
	addRoutes(r, db)

	http.ListenAndServe(":8080", handlers.CORS(getCors()...)(r))
}
