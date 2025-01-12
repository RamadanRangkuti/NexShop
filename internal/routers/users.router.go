package routers

import (
	"github.com/RamadanRangkuti/NexShop/internal/handlers"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func userRouter(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/user")

	var repo repository.UserRepositoryInterface = repository.NewUserRepository(d)
	handler := handlers.NewUserHandler(repo)
	route.GET("/", handler.GetAllUser)
	route.GET("/:id", handler.GetUserById)
	route.POST("/", handler.CreateUser)
	route.PUT("/:id", handler.UpdateUser)
	route.DELETE("/:id", handler.DeleteUser)

	// Get user by username
	route.GET("/username/:username", handler.GetUserByUsername)
	// Get users by signup date
	route.GET("/signup-date", handler.GetUsersBySignupDate)
	//Get user by email
	route.GET("/email/:email", handler.GetUserByEmail)
}
