package handlers

import (
	"fmt"
	"strconv"

	"github.com/RamadanRangkuti/NexShop/internal/dto"
	"github.com/RamadanRangkuti/NexShop/internal/models"
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

	var input dto.InsertCartRequest
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

	// Bind input and validate quantity
	if err := ctx.ShouldBind(&input); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	if input.Quantity <= 0 {
		response.BadRequest("Invalid input", "Quantity must be greater than zero")
		return
	}

	// Validate product availability and stock
	product, err := h.ProductRepositoryInterface.FindProductById(input.ProductID)
	if err != nil || product == nil {
		response.NotFound("Product not found", nil)
		return
	}
	if product.Stock < input.Quantity {
		response.BadRequest("Insufficient stock", "Not enough stock available")
		return
	}

	// Check if the product is already in the cart
	cartItem, err := h.ShoppingCartRepositoryInterface.FindCartItem(id, input.ProductID)
	if err != nil {
		response.InternalServerError("Failed to check cart", err.Error())
		return
	}

	// Prepare data for cart item
	newCart := models.InsertShopping{
		UserID:    id,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}

	var created_cart *models.InsertShopping
	if cartItem != nil {
		// Update existing cart item
		err = h.ShoppingCartRepositoryInterface.UpdateCartItem(id, input.ProductID, input.Quantity)
		if err != nil {
			response.InternalServerError("Failed to update cart item", err.Error())
			return
		}
	} else {
		// Add new item to cart
		created_cart, err = h.ShoppingCartRepositoryInterface.AddCartItem(&newCart)
		if err != nil {
			response.InternalServerError("Failed to create cart item", err.Error())
			return
		}
	}

	var showProduct *models.DetailItemShoppingCart
	if created_cart != nil {
		showProduct, _ = h.ShoppingCartRepositoryInterface.FindCartById(created_cart.ID)
		if err != nil {
			response.InternalServerError("Failed to show cart item", err.Error())
			return
		}
		response.Success("Item added to cart", showProduct)
		return
	}

	updated_cart, _ := h.ShoppingCartRepositoryInterface.FindCartItem(id, input.ProductID)

	response.Success("Item quantity updated ", updated_cart)

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

func (h *ShoppingCartHandler) UpdateCartItem(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	userIdToken, exists := ctx.Get("UserId")

	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	userID, ok := userIdToken.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	var req struct {
		Quantity int `json:"quantity" form:"quantity"`
	}

	if err := ctx.ShouldBind(&req); err != nil {
		response.BadRequest("update user failed", err.Error())
		return
	}

	result, err := h.FindCartById(id)
	if err != nil {
		response.InternalServerError("get cart failed", err.Error())
		return
	}

	if result == nil {
		response.NotFound(fmt.Sprintf("Cart with ID %d not found", id), nil)
		return
	}

	err = h.ShoppingCartRepositoryInterface.EditQuantityCartItem(userID, id, req.Quantity)
	if err != nil {
		response.InternalServerError("Failed to update cart item", err.Error())
		return
	}

	response.Success(fmt.Sprintf("Cart with ID %d update successfully", id), nil)
}

func (h *ShoppingCartHandler) DeleteCartItem(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	result, err := h.FindCartById(id)
	if err != nil {
		response.InternalServerError("get cart failed", err.Error())
		return
	}

	if result == nil {
		response.NotFound(fmt.Sprintf("Cart with ID %d not found", id), nil)
		return
	}

	err = h.RemoveItemCart(id)
	if err != nil {
		response.BadRequest("delete cart failed", err.Error())
		return
	}

	response.Success(fmt.Sprintf("Cart with ID %d deleted successfully", id), nil)
}
