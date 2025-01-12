package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *gin.Engine {
	router := gin.Default()

	authRouter(router, db)
	productRouter(router, db)
	userRouter(router, db)
	accountRouter(router, db)
	cartRouter(router, db)
	orderRouter(router, db)
	return router
}
