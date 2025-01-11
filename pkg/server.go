package pkg

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Server(router *gin.Engine) *http.Server {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	defer fmt.Println("Server successfuly running on PORT :", os.Getenv("PORT"))

	var addr string = "0.0.0.0:8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 10,
		Handler:      router,
	}
	return server
}
