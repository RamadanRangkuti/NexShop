package routers

import (
	"github.com/RamadanRangkuti/NexShop/internal/handlers"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func cartRouter(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/cart")

	var cartRepo repository.ShoppingCartRepositoryInterface = repository.NewShoppingCartRepository(d)
	var productRepo repository.ProductRepositoryInterface = repository.NewProductRepository(d) // Menambahkan ProductRepositoryInterface
	handler := handlers.NewShoppingCartHandler(cartRepo, productRepo)                          // Menyediakan kedua repository
	route.POST("/:id/add", handler.AddToCart)
	route.GET("/:id", handler.GetCartById)
}
