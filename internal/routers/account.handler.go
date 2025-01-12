package routers

import (
	"github.com/RamadanRangkuti/NexShop/internal/handlers"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func accountRouter(g *gin.Engine, d *sqlx.DB) {
	route := g.Group("/account")

	var repo repository.AccountRepositoryInterface = repository.NewAccountRepository(d)
	handler := handlers.NewAccountHandler(repo)
	route.POST("/:id/deposit", handler.Deposit)
	route.POST("/:id/withdraw", handler.Withdraw)
	route.GET("/:id/balance", handler.GetBalance)
}
