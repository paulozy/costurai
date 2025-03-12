package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/paulozy/costurai/internal/infra/database"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

type AuthController struct {
	authenticationUseCase *usecases.AuthenticationUseCase
	sendOtpCodeUseCase    *usecases.SendOTPCodeUseCase
	dressMakerRepository  database.DressmakerRepositoryInterface
}

type AuthUseCasesInput struct {
	AuthUseCase        *usecases.AuthenticationUseCase
	SendOtpCodeUseCase *usecases.SendOTPCodeUseCase
}

func NewAuthController(
	usecases AuthUseCasesInput,
	dr database.DressmakerRepositoryInterface,
) *AuthController {
	return &AuthController{
		authenticationUseCase: usecases.AuthUseCase,
		sendOtpCodeUseCase:    usecases.SendOtpCodeUseCase,
		dressMakerRepository:  dr,
	}
}

func (ac *AuthController) AuthenticateDressmaker(c *gin.Context) {
	var input usecases.AuthenticationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	output, err := ac.authenticationUseCase.DressmakerExecute(input)
	if err.Message != "" {
		c.JSON(err.Status, gin.H{"error": err.Message, "reason": err.Error})
		return
	}

	json := make(map[string]any)

	if output.Dressmaker != nil {
		json["user"] = output.Dressmaker
	} else {
		json["user"] = output.User
	}

	json["token"] = output.Token

	c.JSON(200, json)
}

func (ac *AuthController) AuthenticateUser(c *gin.Context) {
	var input usecases.AuthenticationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	output, err := ac.authenticationUseCase.UserExecute(input)
	if err.Message != "" {
		c.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	json := make(map[string]any)

	if output.User != nil {
		json["user"] = output.User
	} else {
		json["user"] = output.Dressmaker
	}

	json["token"] = output.Token

	c.JSON(200, json)
}

func (ac *AuthController) SendOTPCode(c *gin.Context) {
	ID := c.Param("id")
	LoggedUser := c.GetString("user")

	if ID != LoggedUser {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var input usecases.SendOTPCodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	output, err := ac.sendOtpCodeUseCase.Execute(input)
	if err.Message != "" {
		c.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	fmt.Printf("[SendOTPCodeController]: %v", output)

	c.JSON(200, gin.H{"data": "Code sent succefully"})
}
