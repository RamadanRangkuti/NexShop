package handlers

import (
	"fmt"
	"strings"

	"github.com/RamadanRangkuti/NexShop/internal/dto"
	"github.com/RamadanRangkuti/NexShop/internal/models"
	"github.com/RamadanRangkuti/NexShop/internal/repository"
	"github.com/RamadanRangkuti/NexShop/pkg"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	repository.AuthRepositoryInterface
	repository.UserRepositoryInterface
}

func NewAuthHandler(
	authRepo repository.AuthRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *AuthHandler {
	return &AuthHandler{
		AuthRepositoryInterface: authRepo,
		UserRepositoryInterface: userRepo,
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	var input dto.RegisterRequest

	if err := ctx.ShouldBind(&input); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	if !IsValidEmail(input.Email) {
		response.BadRequest("Invalid email format", nil)
		return
	}

	if !IsValidPassword(input.Password) {
		response.BadRequest("Password must be at least 8 characters long", nil)
		return
	}

	existingUser, err := h.UserRepositoryInterface.FindUserByEmail(input.Email)
	if err != nil {
		response.InternalServerError("Failed to query user", err.Error())
		return
	}

	if existingUser != nil {
		response.BadRequest(fmt.Sprintf("Email %s is already registered", input.Email), nil)
		return
	}
	hashed := pkg.GenerateHash(input.Password)
	fmt.Println("Generated Hash:", hashed)
	newUser := models.Auth{
		Email:    input.Email,
		Username: "",
		Password: hashed,
	}
	createdUser, err := h.AuthRepositoryInterface.RegisterUser(&newUser)
	if err != nil {
		response.InternalServerError("Failed to create user", err.Error())
		return
	}

	response.Created("Success register user", createdUser)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	response := pkg.NewResponse(ctx)
	var input dto.LoginRequest

	if err := ctx.ShouldBind(&input); err != nil {
		response.BadRequest("Invalid input", err.Error())
		return
	}

	user, _ := h.UserRepositoryInterface.FindUserByEmail(input.Email)
	if user == nil || !pkg.VerifyHash(input.Password, user.Password) {
		response.BadRequest("Invalid email or password", nil)
		return
	}

	token, err := pkg.GenerateToken(user.Id)
	if err != nil {
		response.InternalServerError("Failed to generate token", err.Error())
		return
	}

	response.Success("Login success", map[string]interface{}{
		"token": token,
	})
}

func IsValidEmail(email string) bool {
	if len(email) < 5 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}

func IsValidPassword(password string) bool {
	return len(password) >= 8
}
