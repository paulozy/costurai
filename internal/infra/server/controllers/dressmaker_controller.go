package controllers

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paulozy/costurai/internal/infra/database"
	usecases "github.com/paulozy/costurai/internal/usecase"
)

type DressmakerController struct {
	dressMakerRepository    database.DressmakerRepositoryInterface
	dressMakerReviewsRepository database.DressmakerReviewsRepositoryInterface
	createDressmakerUseCase *usecases.CreateDressMakerUseCase
	updateDressmakerUseCase *usecases.UpdateDressMakerUseCase
	getDressmakersByProximityUseCase *usecases.GetDressmakersByProximityUseCase
	getDressmakersByServicesUseCase *usecases.GetDressmakersByServicesUseCase
	addDressmakerReviewUseCase *usecases.AddDressmakerReviewUseCase
}

type DressmakerUseCasesInput struct {
	CreateDressmakerUseCase *usecases.CreateDressMakerUseCase
	UpdateDressmakerUseCase *usecases.UpdateDressMakerUseCase
	GetDressmakersByProximityUseCase *usecases.GetDressmakersByProximityUseCase
	GetDressmakersByServicesUseCase *usecases.GetDressmakersByServicesUseCase
	AddDressmakerReviewUseCase *usecases.AddDressmakerReviewUseCase
}

func NewDressmakerController(dmRepo database.DressmakerRepositoryInterface, dmrRepo database.DressmakerReviewsRepositoryInterface, usecases DressmakerUseCasesInput) *DressmakerController {
	return &DressmakerController{
		dressMakerRepository:    dmRepo,
		dressMakerReviewsRepository: dmrRepo,
		createDressmakerUseCase: usecases.CreateDressmakerUseCase,
		updateDressmakerUseCase: usecases.UpdateDressmakerUseCase,
		getDressmakersByProximityUseCase: usecases.GetDressmakersByProximityUseCase,
		getDressmakersByServicesUseCase: usecases.GetDressmakersByServicesUseCase,
		addDressmakerReviewUseCase: usecases.AddDressmakerReviewUseCase,
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
	searchBy := c.Query("search_by")

	if searchBy == "services"{
		services := c.Query("services")

		if services == "" {
			c.JSON(400, gin.H{"error": "On search by services you must provide the services"})
			return
		}

		var input usecases.GetDressmakersByServicesInput

		input.Services = services

		dressmakers, ucError := dc.getDressmakersByServicesUseCase.Execute(input)
		if ucError.Message != "" {
			c.JSON(ucError.Status, gin.H{"error": ucError.Message, "reason": ucError.Error})
			return
		}

		c.JSON(200, gin.H{"data": dressmakers})

		return
	} else if searchBy == "proximity" {
		println("GetDressmakersByProximity")

		latitude := c.Query("latitude")
		longitude := c.Query("longitude")
		maxDistance := c.Query("max_distance")

		if latitude == "" || longitude == "" {
			c.JSON(400, gin.H{"error": "Latitude and Longitude are required"})
			return
		}

		if maxDistance == "" {
			c.JSON(400, gin.H{"error": "Max distance is required"})
			return
		}

		var input usecases.GetDressmakersByProximityInput

		normalizedLatitude, err := strconv.ParseFloat(latitude, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Latitude must be a float"})
			return
		}

		normalizedLongitude, err := strconv.ParseFloat(longitude, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Longitude must be a float"})
			return
		}

		normalizedMaxDistance, err := strconv.Atoi(maxDistance)
		if err != nil {
			c.JSON(400, gin.H{"error": "Max distance must be an integer"})
			return
		}

		input.Latitude = normalizedLatitude
		input.Longitude = normalizedLongitude
		input.Distance = normalizedMaxDistance

		dressmakers, ucError := dc.getDressmakersByProximityUseCase.Execute(input)
		if ucError.Message != "" {
			c.JSON(ucError.Status, gin.H{"error": ucError.Message, "reason": ucError.Error})
			return
		}

		c.JSON(200, gin.H{"data": dressmakers})

		return
	}
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