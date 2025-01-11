package routers

import (
	"github.com/RamadanRangkuti/NexShop/internal/handlers"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func productRouter(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/product")

	var repo repository.ProductRepositoryInterface = repository.NewProductRepository(d)
	handler := handlers.NewProductHandler(repo)

	route.GET("/", handler.GetAllProduct)
	route.GET("/:id", handler.GetProductById)
	route.POST("/", handler.CreateProduct)
	route.PUT("/:id", handler.UpdateProduct)
	route.DELETE("/:id", handler.DeleteProduct)
}
