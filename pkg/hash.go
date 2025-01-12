package pkg

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pilinux/argon2"
)

func VerifyHash(hash string, password string) bool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	result, err := argon2.ComparePasswordAndHash(hash, os.Getenv("HASHKEY"), password)
	fmt.Println("Error during verification:", err)
	return result
}

func GenerateHash(password string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	result, _ := argon2.CreateHash(password, os.Getenv("HASHKEY"), argon2.DefaultParams)
	return result
}
