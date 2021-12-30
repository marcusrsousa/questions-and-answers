package utils

import (
	"fmt"
	"level7/questions-and-answers/model"
	"level7/questions-and-answers/repository"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("5G$w7vz:>s`SyaU;$P/Y`:9A$v[9a8")

type Claims struct {
	User model.User `json:"user"`
	jwt.StandardClaims
}

func GenerateJWT(user model.User) (string, error) {
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request, model.User), userRepository *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Not Authorized")
			return
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(r.Header["Token"][0], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error on token method")
			}
			return mySigningKey, nil
		})

		if err != nil {
			fmt.Fprint(w, err.Error())
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Not Authorized")
			return
		}

		user := userRepository.FindById(claims.User.ID)
		if user.ID != claims.User.ID {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Not Authorized")
			return
		}

		endpoint(w, r, *user)

	}
}
