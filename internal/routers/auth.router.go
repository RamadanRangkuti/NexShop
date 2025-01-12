package routers

import (
	"github.com/RamadanRangkuti/NexShop/internal/handlers"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func authRouter(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/auth")

	var authRepo repository.AuthRepositoryInterface = repository.NewAuthRepository(d)
	var userRepo repository.UserRepositoryInterface = repository.NewUserRepository(d)
	handler := handlers.NewAuthHandler(authRepo, userRepo)
	route.POST("/register", handler.Register)
	route.POST("/login", handler.Login)
}
