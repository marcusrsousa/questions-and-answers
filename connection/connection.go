package connection

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func GetConnection() *gorm.DB {
	db, err := gorm.Open("postgres", "host=192.168.240.1 port=5432 user=postgres dbname=postgres sslmode=disable password=mysecretpassword")

	if err != nil {
		log.Fatalln("database connection error", err)
	}

	return db
}
