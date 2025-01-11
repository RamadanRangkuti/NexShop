package handlers

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/RamadanRangkuti/NexShop/internal/dto"
	"github.com/RamadanRangkuti/NexShop/internal/models"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/RamadanRangkuti/NexShop/pkg"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	repository.ProductRepositoryInterface
}

func NewProductHandler(r repository.ProductRepositoryInterface) *ProductHandler {
	return &ProductHandler{r}
}

func (h *ProductHandler) GetAllProduct(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))
	search := ctx.DefaultQuery("search", "")
	sort := ctx.DefaultQuery("sort", "id")
	order := ctx.DefaultQuery("order", "asc")

	if page < 1 {
		response.BadRequest("Invalid input", "Page must be 1 or greater")
		return
	}
	if limit < 1 {
		response.BadRequest("Invalid input", "Limit must be 1 or greater")
		return
	}

	if sort != "id" && sort != "price" {
		response.BadRequest("Invalid input", "Sort must be either 'id' or 'price'")
		return
	}
	if order != "asc" && order != "desc" {
		response.BadRequest("Invalid input", "Order must be either 'asc' or 'desc'")
		return
	}

	offset := (page - 1) * limit

	count, _ := h.CountProduct(search)
	products, err := h.FindAllProduct(limit, offset, search, sort, order)
	if err != nil {
		response.InternalServerError("get data failed", err.Error())
		return
	}
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	pageInfo := &pkg.PageInfo{
		CurrentPage: page,
		NextPage:    page + 1,
		PrevPage:    page - 1,
		TotalPage:   totalPages,
		TotalData:   count,
	}
	if page >= totalPages {
		pageInfo.NextPage = 0
	}
	if page <= 1 {
		pageInfo.PrevPage = 0
	}

	response.GetAllSuccess("Success get all products", products, pageInfo)
}

func (h *ProductHandler) GetProductById(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	result, err := h.FindProductById(id)
	if err != nil {
		response.InternalServerError("get data failed", err.Error())
		return
	}

	if result == nil {
		response.NotFound(fmt.Sprintf("Product with ID %d not found", id), nil)
		return
	}
	response.Success("Success get product", result)
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	var req dto.CreateProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadRequest("create data failed", err.Error())
		return
	}

	if err := ValidateProduct(&req.Price, &req.Stock); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	result, err := h.InsertProduct(&product)

	if err != nil {
		response.BadRequest("create data failed", err.Error())
		return
	}
	response.Created("create data success", result)
}

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.BadRequest("Invalid input", "ID must be a valid number")
		return
	}

	existingData, _ := h.FindProductById(id)
	if existingData == nil {
		response.NotFound(fmt.Sprintf("Product with ID %d not found", id), nil)
		return
	}

	var req dto.UpdateProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadRequest("update data failed", err.Error())
		return
	}

	if err := ValidateProduct(req.Price, req.Stock); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	if req.Name != nil {
		existingData.Name = *req.Name
	}
	if req.Description != nil {
		existingData.Description = *req.Description
	}
	if req.Price != nil {
		existingData.Price = *req.Price
	}
	if req.Stock != nil {
		existingData.Stock = *req.Stock
	}

	now := time.Now()
	existingData.Updated_at = &now

	result, err := h.EditProduct(id, existingData)
	if err != nil {
		response.InternalServerError("Failed to update product", err.Error())
		return
	}

	response.Success("Update data success", result)
}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.BadRequest("Invalid input", "ID must be a valid number")
		return
	}
	existingData, _ := h.FindProductById(id)
	if existingData == nil {
		response.NotFound(fmt.Sprintf("Product with ID %d not found", id), nil)
		return
	}
	err = h.RemoveProduct(id)
	if err != nil {
		response.BadRequest("delete data failed", err.Error())
		return
	}

	response.Success(fmt.Sprintf("Product with ID %d deleted successfully", id), nil)
}

func ValidateProduct(price *float64, stock *int) error {
	if price != nil && *price < 0 {
		return errors.New("price cannot be negative")
	}
	if stock != nil && *stock < 0 {
		return errors.New("stock cannot be negative")
	}
	return nil
}
