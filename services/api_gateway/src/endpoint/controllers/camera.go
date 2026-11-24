package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"api_gateway/entities"
	"api_gateway/helpers"
	"api_gateway/services"
)

type CameraController struct {
	cameraService *services.CameraService
}

func NewCameraController(service *services.CameraService) *CameraController {
	return &CameraController{
		cameraService: service,
	}
}

// CreateCamera godoc
// @Summary Create camera with passed properties.
// @Description Creates a new person in the system
// @Param request body entities.Camera true "camera entity" (???)
// @Success 201 {object} int
// @Failure 400
// @Failure 500
// @Router /api/camera/create [post]
// @Security BearerAuth
// @Tags Camera
func (c *CameraController) CreateCamera(ctx *gin.Context) {
	camera := &entities.Camera{}
	err := ctx.ShouldBindJSON(&camera)
	if err != nil {
		log.Printf("Error shouldBindJSON: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	id, err := c.cameraService.CreateCamera(camera)
	if err != nil {
		log.Printf("Error create camera: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, id)
}

// GetCamera godoc
// @Summary Get camera info by ID.
// @Produce json
// @Param id path integer true "Camera ID"
// @Success 200 {object} entities.Camera
// @Failure 400
// @Failure 404
// @Failure 500 "Internal server error"
// @Router /api/camera/{id} [get]
// @Security BearerAuth
// @Tags Camera
func (c *CameraController) GetCameraByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Error parse id: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	camera, err := c.cameraService.GetCameraByID(id)
	if err != nil {
		log.Printf("Error get camera: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, camera)
}

// GetCameraList godoc
// @Summary Get paginated list of cameras
// @Produce json
// @Description Retrieves a paginated list of cameras
// @Param limit query integer false "limit"
// @Param offset query integer false "offset"
// @Success 200 {array} entities.Camera "List of cameras"
// @Failure 400 "Invalid pagination parameters"
// @Failure 500 "Internal server error"
// @Router /api/camera/list [get]
// @Security BearerAuth
// @Tags Camera
func (c *CameraController) GetCameraList(ctx *gin.Context) {
	limit, offset, err := helpers.ExtractPagination(ctx.Request)
	if err != nil {
		log.Printf("Invalid pagination. Set to default")
	}

	cameras, err := c.cameraService.GetCameraList(int(limit), int(offset))
	if err != nil {
		log.Printf("Error get camera list: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, cameras)
}

// GetCount godoc
// @Summary Get cameras count
// @Description Returns the total number of cameras in the system
// @Produce json
// @Param host query string false "host"
// @Param port query string false "port"
// @Success 200 {integer} integer "Number of cameras"
// @Failure 500 "Internal server error"
// @Router /api/camera/count [get]
// @Security BearerAuth
// @Tags Camera
func (c *CameraController) GetCameraCount(ctx *gin.Context) {
	count, err := c.cameraService.GetCameraCount()
	if err != nil {
		log.Printf("Error get camera count: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, count)
}

// UpdateCamera godoc
// @Summary Update camera by ID.
// @Description Updates information for an existing person record
// @Accept  json
// @Produce json
// @Param id path integer true "Camera ID"
// @Param request body entities.Camera true "camera entity" (???)
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500 "Internal server error"
// @Router /api/camera/update/{id} [put]
// @Security BearerAuth
// @Tags Camera
func (c *CameraController) UpdateCamera(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Error parse id: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	camera := &entities.Camera{}
	err = ctx.ShouldBindJSON(&camera)
	if err != nil {
		log.Printf("Error shouldBindJSON: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = c.cameraService.UpdateCamera(id, camera)
	if err != nil {
		log.Printf("Error update camera: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteCamera godoc
// @Summary Delete camera object by ID.
// @Accept  json
// @Produce json
// @Param id path integer true "Camera ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500 "Internal server error"
// @Router /api/camera/{id} [delete]
// @Security BearerAuth
// @Tags Camera
func (c *CameraController) DeleteCamera(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Error parse id: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = c.cameraService.DeleteCamera(id)
	if err != nil {
		log.Printf("Error delete camera: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
