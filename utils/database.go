package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func GetConnection() *sqlx.DB {
	var portString string = goDotEnvVariable("DB_PORT")
	port, err := strconv.Atoi(portString)

	if err != nil {
		log.Fatalf("Error loading port from .env file")
	  }

	
	
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"postgres", 5432, 
		"chetbackiewicz", "i90jsadfp9030jf!!", 
		"elo_sport")

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	log.Println("DB Connection established...")
	return db
}

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

