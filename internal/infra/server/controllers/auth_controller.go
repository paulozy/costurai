package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/paulozy/costurai/internal/infra/database"
	usecases "github.com/paulozy/costurai/internal/usecase/auth"
	"github.com/paulozy/costurai/internal/usecase/auth/dtos"
)

type AuthController struct {
	authDressmakerUseCase *usecases.AuthDressmakerUseCase
	authUserUseCase       *usecases.AuthUserUseCase
	sendOTPUseCase        *usecases.SendOTPUseCase
	dressmakerRepository  database.DressmakerRepositoryInterface
	userRepository        database.UserRepositoryInterface
}

func NewAuthController(
	authDressmakerUseCase *usecases.AuthDressmakerUseCase,
	authUserUseCase *usecases.AuthUserUseCase,
	sendOTPUseCase *usecases.SendOTPUseCase,
	dr database.DressmakerRepositoryInterface,
	ur database.UserRepositoryInterface,
) *AuthController {
	return &AuthController{
		authDressmakerUseCase: authDressmakerUseCase,
		authUserUseCase:       authUserUseCase,
		sendOTPUseCase:        sendOTPUseCase,
		dressmakerRepository:  dr,
		userRepository:        ur,
	}
}

func (ac *AuthController) AuthenticateDressmaker(c *gin.Context) {
	var input dtos.AuthenticationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := ac.authDressmakerUseCase.Execute(input)
	if err.Message != "" {
		c.JSON(err.Status, gin.H{"error": err.Message, "reason": err.Error})
		return
	}

	c.JSON(200, gin.H{"data": token})
}

func (ac *AuthController) AuthenticateUser(c *gin.Context) {
	var input dtos.AuthenticationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := ac.authDressmakerUseCase.Execute(input)
	if err.Message != "" {
		c.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (ac *AuthController) SendOTP(c *gin.Context) {
	var input dtos.OTPSendAndVerifyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := ac.sendOTPUseCase.Execute(input)
	if err.Message != "" {
		c.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	c.JSON(200, gin.H{"data": "Code sent succefully"})
}
