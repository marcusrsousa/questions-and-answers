package main

import (
	"level7/questions-and-answers/connection"
	"level7/questions-and-answers/controller"
	"level7/questions-and-answers/model"
	"level7/questions-and-answers/repository"
	"level7/questions-and-answers/utils"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func addRoutes(r *mux.Router, db *gorm.DB) {

	qc, lc, ac := getControllers(db)

	r.HandleFunc("/question", utils.IsAuthorized(qc.Add, lc.Repository)).Methods(http.MethodPost)
	r.HandleFunc("/question/{id}", utils.IsAuthorized(qc.Update, lc.Repository)).Methods(http.MethodPut)
	r.HandleFunc("/question/{id}", utils.IsAuthorized(qc.Delete, lc.Repository)).Methods(http.MethodDelete)
	r.HandleFunc("/question/{id}", utils.IsAuthorized(qc.GetById, lc.Repository)).Methods(http.MethodGet)
	r.HandleFunc("/question", qc.Get).Methods(http.MethodGet)
	r.HandleFunc("/question/{questionId}/answer", ac.Get).Methods(http.MethodGet)
	r.HandleFunc("/question/{questionId}/answer", utils.IsAuthorized(ac.Add, lc.Repository)).Methods(http.MethodPost)
	r.HandleFunc("/question/{questionId}/answer/{id}", utils.IsAuthorized(ac.Update, lc.Repository)).Methods(http.MethodPut)
	r.HandleFunc("/question/{questionId}/answer/{id}", utils.IsAuthorized(ac.Delete, lc.Repository)).Methods(http.MethodDelete)
	r.HandleFunc("/login", lc.Login).Methods(http.MethodPost)
	r.HandleFunc("/signup", lc.SignUp).Methods(http.MethodPost)
	r.HandleFunc("/change-password", utils.IsAuthorized(lc.ChangePassword, lc.Repository)).Methods(http.MethodPost)

}

func getControllers(db *gorm.DB) (*controller.QuestionController, *controller.LoginController, *controller.AnswerController) {
	questionsRep := &repository.QuestionRepository{Db: db}
	answersRep := &repository.AnswerRepository{Db: db}
	usersRep := &repository.UserRepository{Db: db}
	return &controller.QuestionController{Repository: questionsRep}, &controller.LoginController{Repository: usersRep}, &controller.AnswerController{Repository: answersRep}
}

func initialMigration(db *gorm.DB) {
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Question{})
	db.AutoMigrate(&model.Answer{})
}

func getCors() []handlers.CORSOption {
	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	return []handlers.CORSOption{header, methods, origins}
}

func main() {
	db := connection.GetConnection()
	db.LogMode(true)
	initialMigration(db)
	r := mux.NewRouter()
	addRoutes(r, db)

	http.ListenAndServe(":8080", handlers.CORS(getCors()...)(r))
}
