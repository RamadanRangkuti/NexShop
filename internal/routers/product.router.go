package routers

import (
	"github.com/RamadanRangkuti/NexShop/internal/handlers"
	"github.com/RamadanRangkuti/NexShop/internal/middlewares"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func productRouter(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/product")

	var repo repository.ProductRepositoryInterface = repository.NewProductRepository(d)
	handler := handlers.NewProductHandler(repo)

	route.GET("/", handler.GetAllProduct)
	route.GET("/:id", middlewares.ValidateToken(), handler.GetProductById)
	route.POST("/", middlewares.ValidateToken(), handler.CreateProduct)
	route.PUT("/:id", middlewares.ValidateToken(), handler.UpdateProduct)
	route.DELETE("/:id", middlewares.ValidateToken(), handler.DeleteProduct)
}
