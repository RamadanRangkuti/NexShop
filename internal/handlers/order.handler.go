package handlers

import (
	"fmt"

	"github.com/RamadanRangkuti/NexShop/internal/dto"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/RamadanRangkuti/NexShop/pkg"
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

func (h *PurchaseHandler) CompletePurchase(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	userID, exists := ctx.Get("UserId")
	if !exists {
		response.Unauthorized("Unauthorized: User ID not found in token", nil)
		return
	}

	id, ok := userID.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	fmt.Printf("Debug: UserID = %d\n", id)

	// Step 2: Get cart items

	cartItems, err := h.ShoppingCartRepo.FindCartByUserid(id)
	if err != nil {
		fmt.Printf("Error fetching cart items: %v\n", err)
		response.InternalServerError("Failed to fetch shopping cart", err.Error())
		return
	}

	if len(*cartItems) == 0 {
		response.BadRequest("Shopping cart is empty", nil)
		return
	}

	fmt.Printf("Debug: CartItems = %+v\n", *cartItems)

	// Step 3: Validate stock and calculate total
	var total float64
	for _, item := range *cartItems {
		fmt.Println(item.ID)
		product, err := h.ProductRepo.FindProductById(item.ProductID)
		if err != nil {
			fmt.Printf("Error fetching product: %v\n", err)
			response.InternalServerError("Failed to fetch product details", err.Error())
			return
		}
		if product.Stock < item.Quantity {
			response.BadRequest("Insufficient stock for product: "+item.ProductName, nil)
			return
		}
		total += product.Price * float64(item.Quantity)
	}

	fmt.Printf("Debug: Total = %.2f\n", total)

	// // Step 4: Validate user balance
	balance, err := h.AccountRepo.GetBalance(id)
	if err != nil {
		fmt.Printf("Error fetching user balance: %v\n", err)
		response.InternalServerError("Failed to fetch user balance", err.Error())
		return
	}
	if balance < total {
		response.BadRequest("Insufficient balance", nil)
		return
	}

	fmt.Printf("Debug: Balance = %.2f\n", balance)

	// // Step 5: Create order
	order, err := h.OrderRepo.CreateOrder(id, total)
	if err != nil {
		fmt.Printf("Error creating order: %v\n", err)
		response.InternalServerError("Failed to create order", err.Error())
		return
	}

	fmt.Printf("Debug: Order created with ID = %d\n", order.ID)

	// // Step 6: Add order items
	for _, item := range *cartItems {
		product, _ := h.ProductRepo.FindProductById(item.ProductID)
		err := h.OrderRepo.AddOrderItem(order.ID, product.Id, item.Quantity, product.Price)
		if err != nil {
			fmt.Printf("Error adding order item: %v\n", err)
			response.InternalServerError("Failed to add order item", err.Error())
			return
		}
	}

	// // Step 7: Reduce stock
	for _, item := range *cartItems {
		fmt.Println(item)
		err := h.ProductRepo.ReduceStock(item.ProductID, item.Quantity)
		if err != nil {
			fmt.Printf("Error reducing stock: %v\n", err)
			response.InternalServerError("Failed to reduce stock", err.Error())
			return
		}
	}

	// // Step 8: Deduct balance
	err = h.AccountRepo.Withdraw(id, total)
	if err != nil {
		fmt.Printf("Error deducting balance: %v\n", err)
		response.InternalServerError("Failed to deduct balance", err.Error())
		return
	}

	// // Step 9: Clear shopping cart
	err = h.ShoppingCartRepo.ClearCartByUserID(id)
	if err != nil {
		fmt.Printf("Error clearing shopping cart: %v\n", err)
		response.InternalServerError("Failed to clear cart", err.Error())
		return
	}

	// // Step 10: Return successful response
	response.Success("Purchase completed successfully", dto.CompletePurchaseResponse{
		OrderID:    order.ID,
		TotalPrice: total,
		Message:    "Purchase completed successfully",
	})
}
