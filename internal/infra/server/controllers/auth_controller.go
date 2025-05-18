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
	VerifyOTPUseCase      *usecases.VerifyOTPUseCase
	dressmakerRepository  database.DressmakerRepositoryInterface
	userRepository        database.UserRepositoryInterface
}

func NewAuthController(
	authDressmakerUseCase *usecases.AuthDressmakerUseCase,
	authUserUseCase *usecases.AuthUserUseCase,
	sendOTPUseCase *usecases.SendOTPUseCase,
	verifyOTPUseCase *usecases.VerifyOTPUseCase,
	dr database.DressmakerRepositoryInterface,
	ur database.UserRepositoryInterface,
) *AuthController {
	return &AuthController{
		authDressmakerUseCase: authDressmakerUseCase,
		authUserUseCase:       authUserUseCase,
		sendOTPUseCase:        sendOTPUseCase,
		VerifyOTPUseCase:      verifyOTPUseCase,
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
	var input dtos.SendOTPInput
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

func (ac *AuthController) VerifyOTP(c *gin.Context) {
	var input dtos.VerifyOTPInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	LoggedUser := c.GetString("user")

	switch c.FullPath() {
	case "/otp/dressmaker/verify":
		input.Enabling = "dressmaker"
		input.DressmakerID = LoggedUser
	case "/otp/user/verify":
		input.Enabling = "user"
		input.UserID = LoggedUser
	default:
		input.Enabling = ""
	}

	err := ac.VerifyOTPUseCase.Execute(input)
	if err.Message != "" {
		c.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	c.JSON(200, gin.H{"data": "OTP verified successfully"})
}
