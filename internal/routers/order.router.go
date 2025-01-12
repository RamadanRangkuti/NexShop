package routers

import (
	"github.com/RamadanRangkuti/NexShop/internal/handlers"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func orderRouter(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/order")

	cartRepo := repository.NewShoppingCartRepository(d)
	productRepo := repository.NewProductRepository(d)
	accountRepo := repository.NewAccountRepository(d)
	orderRepo := repository.NewOrderRepository(d)

	handler := handlers.NewPurchaseHandler(cartRepo, productRepo, accountRepo, orderRepo)

	route.POST("/purchase", handler.CompletePurchase)
}
