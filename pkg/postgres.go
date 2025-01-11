package pkg

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func PostgresConnection() (*sqlx.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	log.Printf("DB_HOST=%s, DB_USER=%s, DB_PASS=%s, DB_NAME=%s", host, user, pass, dbname)

	if host == "" || user == "" || pass == "" || dbname == "" {
		return nil, fmt.Errorf("database configuration is incomplete")
	}

	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, pass, dbname)

	return sqlx.Connect("postgres", config)
}
