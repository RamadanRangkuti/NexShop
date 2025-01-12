package handlers

import (
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/RamadanRangkuti/NexShop/pkg"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	repository.AccountRepositoryInterface
}

func NewAccountHandler(r repository.AccountRepositoryInterface) *AccountHandler {
	return &AccountHandler{r}
}

func (h *AccountHandler) Deposit(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	userId, exists := ctx.Get("UserId")

	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)

	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	var req struct {
		Amount float64 `json:"amount" form:"amount"`
	}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	if req.Amount <= 0 {
		response.BadRequest("Invalid input", "Amount must be greater than zero")
		return
	}

	err := h.AccountRepositoryInterface.Deposit(id, req.Amount)
	if err != nil {
		response.InternalServerError("Failed to deposit", err.Error())
		return
	}

	response.Success("Deposit successful", nil)
}

func (h *AccountHandler) Withdraw(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	userId, exists := ctx.Get("UserId")

	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	var req struct {
		Amount float64 `json:"amount" form:"amount"`
	}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	if req.Amount <= 0 {
		response.BadRequest("Invalid input", "Amount must be greater than zero")
		return
	}

	err := h.AccountRepositoryInterface.Withdraw(id, req.Amount)
	if err != nil {
		response.InternalServerError("Failed to withdraw", err.Error())
		return
	}

	response.Success("Withdraw success", nil)
}

func (h *AccountHandler) GetBalance(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	userId, exists := ctx.Get("UserId")

	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	balance, err := h.AccountRepositoryInterface.GetBalance(id)
	if err != nil {
		response.InternalServerError("Failed to fetch balance", err.Error())
		return
	}

	response.Success("Success get Balance", map[string]interface{}{
		"balance": balance,
	})
}
