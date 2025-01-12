package handlers

import (
	"fmt"

	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/RamadanRangkuti/NexShop/pkg"
	"github.com/gin-gonic/gin"
)

type ShoppingCartHandler struct {
	repository.ShoppingCartRepositoryInterface
	repository.ProductRepositoryInterface
}

func NewShoppingCartHandler(
	cartRepo repository.ShoppingCartRepositoryInterface,
	productRepo repository.ProductRepositoryInterface,
) *ShoppingCartHandler {
	return &ShoppingCartHandler{
		ShoppingCartRepositoryInterface: cartRepo,
		ProductRepositoryInterface:      productRepo,
	}
}

func (h *ShoppingCartHandler) AddToCart(ctx *gin.Context) {
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
		ProductID int `json:"productId" form:"product_id"`
		Quantity  int `json:"quantity" form:"quantity"`
	}
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	// Validasi quntity
	if req.Quantity <= 0 {
		response.BadRequest("Invalid input", "Quantity must be greater than zero")
		return
	}

	// Check apakah product ada dan cukup stock
	product, err := h.ProductRepositoryInterface.FindProductById(req.ProductID)
	if err != nil || product == nil {
		response.NotFound("Product not found", nil)
		return
	}
	if product.Stock < req.Quantity {
		response.BadRequest("Insufficient stock", "Not enough stock available")
		return
	}

	// Check apakah product sudah ada di cart
	cartItem, err := h.ShoppingCartRepositoryInterface.FindCartItem(id, req.ProductID)
	if err != nil {
		response.InternalServerError("Failed to check cart", err.Error())
		return
	}

	if cartItem != nil {
		// Update item di cart yang ada
		err = h.ShoppingCartRepositoryInterface.UpdateCartItem(id, req.ProductID, req.Quantity)
	} else {
		// Tambah item baru ke cart
		err = h.ShoppingCartRepositoryInterface.AddCartItem(id, req.ProductID, req.Quantity)
	}

	if err != nil {
		response.InternalServerError("Failed to update cart", err.Error())
		return
	}

	response.Success("Item added to cart", nil)
}

func (h *ShoppingCartHandler) GetCartById(ctx *gin.Context) {
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

	result, err := h.FindCartByUserid(id)
	if err != nil {
		response.InternalServerError("get cart failed", err.Error())
		return
	}

	//Seharusnya gak dipakai jika menggunakan token
	if result == nil {
		response.NotFound(fmt.Sprintf("Cart with ID %d not found", id), nil)
		return
	}
	response.Success("Success get cart", result)
}
