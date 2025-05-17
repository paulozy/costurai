package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulozy/costurai/internal/infra/database"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

type DressmakerController struct {
	dressMakerRepository             database.DressmakerRepositoryInterface
	dressMakerReviewsRepository      database.DressmakerReviewsRepositoryInterface
	createDressmakerUseCase          *usecases.CreateDressMakerUseCase
	updateDressmakerUseCase          *usecases.UpdateDressMakerUseCase
	getDressmakersByProximityUseCase *usecases.GetDressmakersByProximityUseCase
	addDressmakerReviewUseCase       *usecases.AddDressmakerReviewUseCase
	showDressmakerUseCase            *usecases.ShowDressMakerUseCase
}

type DressmakerUseCasesInput struct {
	CreateDressmakerUseCase          *usecases.CreateDressMakerUseCase
	UpdateDressmakerUseCase          *usecases.UpdateDressMakerUseCase
	GetDressmakersByProximityUseCase *usecases.GetDressmakersByProximityUseCase
	AddDressmakerReviewUseCase       *usecases.AddDressmakerReviewUseCase
	ShowDressmakerUseCase            *usecases.ShowDressMakerUseCase
}

func NewDressmakerController(dmRepo database.DressmakerRepositoryInterface, dmrRepo database.DressmakerReviewsRepositoryInterface, usecases DressmakerUseCasesInput) *DressmakerController {
	return &DressmakerController{
		dressMakerRepository:             dmRepo,
		dressMakerReviewsRepository:      dmrRepo,
		createDressmakerUseCase:          usecases.CreateDressmakerUseCase,
		updateDressmakerUseCase:          usecases.UpdateDressmakerUseCase,
		getDressmakersByProximityUseCase: usecases.GetDressmakersByProximityUseCase,
		addDressmakerReviewUseCase:       usecases.AddDressmakerReviewUseCase,
		showDressmakerUseCase:            usecases.ShowDressmakerUseCase,
	}
}

func (dc *DressmakerController) CreateDressmaker(c *gin.Context) {
	var input usecases.CreateDressMakerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dressmaker, err := dc.createDressmakerUseCase.Execute(input)
	if err.Message != "" {
		c.JSON(err.Status, gin.H{"error": err.Message, "reason": err.Error})
		return
	}

	c.JSON(201, gin.H{"data": dressmaker})
}

func (dc *DressmakerController) UpdateDressmaker(c *gin.Context) {
	ID := c.Param("id")
	LoggedUser := c.GetString("user")

	if ID != LoggedUser {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var input usecases.UpdateDressMakerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	input.ID = ID

	dressmaker, err := dc.updateDressmakerUseCase.Execute(input)
	if err.Message != "" {
		log.Println(err)
		c.JSON(err.Status, gin.H{"error": err.Message})
		return
	}

	c.JSON(200, gin.H{"data": dressmaker})
}

func (dc *DressmakerController) GetDressmakers(c *gin.Context) {
	var input usecases.GetDressmakersByProximityInput

	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Limit == 0 {
		input.Limit = 10
	}

	if input.Page == 0 {
		input.Page = 1
	}

	dressmakers, ucError := dc.getDressmakersByProximityUseCase.Execute(input)
	if ucError.Message != "" {
		c.JSON(ucError.Status, gin.H{"error": ucError.Message, "reason": ucError.Error})
		return
	}

	c.JSON(200, gin.H{"items": dressmakers.Items, "pagination": dressmakers.PaginationInfo})
}

func (dc *DressmakerController) GetDressmaker(c *gin.Context) {
	var input usecases.ShowDressMakerInput

	ID := c.Param("id")
	input.ID = ID

	dressmaker, ucError := dc.showDressmakerUseCase.Execute(input)

	if ucError.Message != "" {
		c.JSON(ucError.Status, gin.H{"error": ucError.Message, "reason": ucError.Error})
		return
	}
	c.JSON(200, gin.H{"data": dressmaker})
}

func (dc *DressmakerController) AddDressmakerReview(c *gin.Context) {
	userID := c.GetString("user")

	var input usecases.AddDressmakerReviewUseCaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error(), "reason": err})
		return
	}

	input.UserID = userID

	dressmaker, ucError := dc.addDressmakerReviewUseCase.Execute(input)
	if ucError.Message != "" {
		c.JSON(ucError.Status, gin.H{"error": ucError.Message, "reason": ucError.Error})
		return
	}

	c.JSON(200, gin.H{"data": dressmaker})
}
