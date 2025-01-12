package routers

import (
	"github.com/RamadanRangkuti/NexShop/internal/handlers"
	"github.com/RamadanRangkuti/NexShop/internal/middlewares"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func userRouter(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/user")

	var repo repository.UserRepositoryInterface = repository.NewUserRepository(d)
	handler := handlers.NewUserHandler(repo)
	route.GET("/", handler.GetAllUser)
	route.GET("/:id", middlewares.ValidateToken(), handler.GetUserById)
	route.POST("/", middlewares.ValidateToken(), handler.CreateUser)
	route.PUT("/:id", middlewares.ValidateToken(), handler.UpdateUser)
	route.DELETE("/:id", middlewares.ValidateToken(), handler.DeleteUser)

	// Get user by username
	route.GET("/username/:username", handler.GetUserByUsername)
	// Get users by signup date
	route.GET("/signup-date", middlewares.ValidateToken(), handler.GetUsersBySignupDate)
	//Get user by email
	route.GET("/email/:email", middlewares.ValidateToken(), handler.GetUserByEmail)
}
