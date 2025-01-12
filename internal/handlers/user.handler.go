package handlers

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/RamadanRangkuti/NexShop/internal/dto"
	"github.com/RamadanRangkuti/NexShop/internal/models"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/RamadanRangkuti/NexShop/pkg"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repository.UserRepositoryInterface
}

func NewUserHandler(r repository.UserRepositoryInterface) *UserHandler {
	return &UserHandler{r}
}

func (h *UserHandler) GetAllUser(ctx *gin.Context) {
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

	count, _ := h.CountUser(search)
	products, err := h.FindAllUser(limit, offset, search, sort, order)
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

	response.GetAllSuccess("Success get all users", products, pageInfo)
}

func (h *UserHandler) GetUserById(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	userId, exists := ctx.Get("UserId")
	fmt.Println("userid", userId)
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)
	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	result, err := h.FindUserById(id)
	if err != nil {
		response.InternalServerError("get user failed", err.Error())
		return
	}

	if result == nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}
	response.Success("Success get user", result)
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	var req dto.CreateUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadRequest("create user failed", err.Error())
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	result, err := h.InsertUser(&user)

	if err != nil {
		response.BadRequest("create user failed", err.Error())
		return
	}
	response.Created("create user success", result)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)

	userId, exists := ctx.Get("UserId")
	fmt.Println("userid", userId)
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)

	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}

	existingData, _ := h.FindUserById(id)
	if existingData == nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadRequest("update user failed", err.Error())
		return
	}

	if req.Email != nil {
		existingData.Email = *req.Email
	}
	if req.Username != nil {
		existingData.Username = *req.Username
	}
	if req.Password != nil {
		if !IsValidPassword(*req.Password) {
			response.BadRequest("Password must be at least 8 characters long", nil)
			return
		}
		existingData.Password = pkg.GenerateHash(*req.Password)
	}

	now := time.Now()
	existingData.Updated_at = &now

	result, err := h.EditUser(id, existingData)
	if err != nil {
		response.InternalServerError("Failed to update user", err.Error())
		return
	}

	response.Success("Update user success", result)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	userId, exists := ctx.Get("UserId")
	fmt.Println("userid", userId)
	if !exists {
		response.Unauthorized("Unauthorized", nil)
		return
	}
	id, ok := userId.(int)

	if !ok {
		response.InternalServerError("Failed to parse user ID from token", nil)
		return
	}
	existingData, _ := h.FindUserById(id)
	if existingData == nil {
		response.NotFound(fmt.Sprintf("User with ID %d not found", id), nil)
		return
	}
	err := h.RemoveUser(id)
	if err != nil {
		response.BadRequest("delete user failed", err.Error())
		return
	}

	response.Success(fmt.Sprintf("User with ID %d deleted successfully", id), nil)
}

func (h *UserHandler) GetUserByUsername(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	username := ctx.Param("username")

	result, err := h.FindUserByUsername(username)
	if err != nil {
		response.InternalServerError("Failed to fetch user by username", err.Error())
		return
	}

	if result == nil {
		response.NotFound(fmt.Sprintf("User with username %s not found", username), nil)
		return
	}

	response.Success("Success get user by username", result)
}

func (h *UserHandler) GetUsersBySignupDate(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	date := ctx.DefaultQuery("date", "2023-01-01")

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		response.BadRequest("Invalid date format", "Date must be in format YYYY-MM-DD")
		return
	}

	result, err := h.FindUsersBySignupDate(date)
	if err != nil {
		response.InternalServerError("Failed to fetch users by signup date", err.Error())
		return
	}

	response.Success(fmt.Sprintf("Success get users who signed up after %s", date), result)
}

func (h *UserHandler) GetUserByEmail(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	email := ctx.Param("email")

	result, err := h.FindUserByEmail(email)
	if err != nil {
		response.InternalServerError("Failed to fetch user by email", err.Error())
		return
	}

	if result == nil {
		response.NotFound(fmt.Sprintf("User with email %s not found", email), nil)
		return
	}

	response.Success("Success get user by email", result)
}

// Fungsi untuk menyaring karakter berbahaya
func SanitizeString(input string) string {
	// Contoh sederhana untuk menghapus karakter-karakter tertentu
	return strings.ReplaceAll(input, "'", "")
}
