package main

import (
	"log"

	"github.com/RamadanRangkuti/NexShop/internal/routers"
	"github.com/RamadanRangkuti/NexShop/pkg"
)

func main() {
	db, err := pkg.PostgresConnection()
	if err != nil {
		log.Fatal(err)
	}
	router := routers.New(db)
	server := pkg.Server(router)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
