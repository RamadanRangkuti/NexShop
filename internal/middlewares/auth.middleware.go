package middlewares

import (
	"fmt"
	"strings"

	"github.com/RamadanRangkuti/NexShop/pkg"
	"github.com/gin-gonic/gin"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := pkg.NewResponse(c)

		head := c.GetHeader("Authorization")
		fmt.Println("token dari header", head)
		if head == "" {
			response.Unauthorized("Unauthorized", nil)
			return
		}
		token := strings.Split(head, " ")[1]
		fmt.Println("token jadi", token)

		if token == "" {
			response.Unauthorized("Unauthorized", nil)
		}

		claims, err := pkg.VerifyToken(token)
		if err != nil {
			response.Unauthorized("Unauthorized", nil)
		}
		c.Set("UserId", claims.UserId)
		c.Next()
	}
}
