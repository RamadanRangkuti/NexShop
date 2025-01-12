package handlers

import (
	"fmt"
	"net/http"

	"github.com/RamadanRangkuti/NexShop/internal/dto"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	ShoppingCartRepo repository.ShoppingCartRepositoryInterface
	ProductRepo      repository.ProductRepositoryInterface
	AccountRepo      repository.AccountRepositoryInterface
	OrderRepo        repository.OrderRepositoryInterface
}

func NewPurchaseHandler(
	cartRepo repository.ShoppingCartRepositoryInterface,
	productRepo repository.ProductRepositoryInterface,
	accountRepo repository.AccountRepositoryInterface,
	orderRepo repository.OrderRepositoryInterface,
) *PurchaseHandler {
	return &PurchaseHandler{
		ShoppingCartRepo: cartRepo,
		ProductRepo:      productRepo,
		AccountRepo:      accountRepo,
		OrderRepo:        orderRepo,
	}
}

func (h *PurchaseHandler) CompletePurchase(c *gin.Context) {
	var req dto.CompletePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 1: Get cart items
	cartItems, err := h.ShoppingCartRepo.FindCartByUserid(req.UserID)
	if err != nil || len(*cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Shopping cart is empty"})
		return
	}

	// Step 2: Validate stock and calculate total
	var total float64
	for _, item := range *cartItems {
		product, err := h.ProductRepo.FindProductById(item.ID)
		if err != nil || product.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product: " + item.ProductName})
			return
		}
		total += product.Price * float64(item.Quantity)
	}

	// Step 3: Validate user balance
	balance, err := h.AccountRepo.GetBalance(req.UserID)
	if err != nil || balance < total {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	fmt.Println(req.UserID)
	// Step 4: Create order and order items
	order, err := h.OrderRepo.CreateOrder(req.UserID, total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	for _, item := range *cartItems {
		product, _ := h.ProductRepo.FindProductById(item.ID)
		err := h.OrderRepo.AddOrderItem(order.ID, product.Id, item.Quantity, product.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add order item"})
			return
		}
	}

	// Step 5: Reduce stock
	for _, item := range *cartItems {
		err := h.ProductRepo.ReduceStock(item.ID, item.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reduce stock"})
			return
		}
	}

	// Step 6: Deduct balance
	err = h.AccountRepo.Withdraw(req.UserID, total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deduct balance"})
		return
	}

	// Step 7: Clear shopping cart
	err = h.ShoppingCartRepo.ClearCartByUserID(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, dto.CompletePurchaseResponse{
		OrderID:    order.ID,
		TotalPrice: total,
		Message:    "Purchase completed successfully",
	})
}
