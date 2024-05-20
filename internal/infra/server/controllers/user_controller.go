package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/paulozy/costurai/internal/infra/database"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

type UserController struct {
	userRepository database.UserRepositoryInterface
	createUserUseCase *usecases.CreateUserUseCase
}

type UserUseCasesInput struct {
	CreateUserUseCase *usecases.CreateUserUseCase
}

func NewUserController(ur database.UserRepositoryInterface, usecases UserUseCasesInput) *UserController {
	return &UserController{
		userRepository: ur,
		createUserUseCase: usecases.CreateUserUseCase,
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var input usecases.CreateUserUseCaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.createUserUseCase.Execute(input)
	if err.Message != "" {
		c.JSON(err.Status, gin.H{"error": err.Message, "reason": err.Error})
		return
	}

	c.JSON(201, gin.H{"data": user})
}