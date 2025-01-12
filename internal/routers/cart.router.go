package routers

import (
	"github.com/RamadanRangkuti/NexShop/internal/handlers"
	"github.com/RamadanRangkuti/NexShop/internal/middlewares"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func cartRouter(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/cart")

	var cartRepo repository.ShoppingCartRepositoryInterface = repository.NewShoppingCartRepository(d)
	var productRepo repository.ProductRepositoryInterface = repository.NewProductRepository(d)
	handler := handlers.NewShoppingCartHandler(cartRepo, productRepo)
	route.POST("/:id/add", middlewares.ValidateToken(), handler.AddToCart)
	route.GET("/:id", middlewares.ValidateToken(), handler.GetCartById)
}
