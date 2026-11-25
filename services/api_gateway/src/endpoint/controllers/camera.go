package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

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
// @Description Creates a new camera in the system
// @Param request body entities.Camera true "camera entity"
// @Success 201 {object} int
// @Failure 400
// @Failure 500
// @Router /api/camera [post]
// @Security BearerAuth
// @Tags Camera
func (c *CameraController) CreateCamera(ctx *gin.Context) {
	camera := &entities.Camera{}
	err := ctx.ShouldBindJSON(camera)
	if err != nil || strings.TrimSpace(camera.IP) == "" || camera.Port < 1 || camera.Port > 65535 {
		log.Printf("Error shouldBindJSON: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	id, err := c.cameraService.CreateCamera(camera)
	if err != nil {
		log.Printf("Error create camera: %v", err)
		writeCameraError(ctx, err)
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
		writeCameraError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, camera)
}

// GetCameraList godoc
// @Summary Get paginated list of cameras
// @Produce json
// @Param camera_id query integer false "Camera ID"
// @Param ip query string false "Camera IP"
// @Param is_enabled query boolean false "Camera enabled state"
// @Param created_from query string false "Inclusive creation time lower bound"
// @Param created_to query string false "Inclusive creation time upper bound"
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
	filter, err := helpers.ExtractFilterFromQueryParams(ctx.Request.URL.Query())
	if err != nil {
		log.Printf("Invalid filter: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	limit, offset, err := helpers.ExtractPagination(ctx.Request)
	if err != nil {
		log.Printf("Invalid pagination. Set to default")
	}

	cameras, err := c.cameraService.GetCameraList(int(limit), int(offset), filter)
	if err != nil {
		log.Printf("Error get camera list: %v", err)
		writeCameraError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, cameras)
}

// GetCount godoc
// @Summary Get cameras count
// @Description Returns the total number of cameras in the system
// @Produce json
// @Param camera_id query integer false "Camera ID"
// @Param ip query string false "Camera IP"
// @Param is_enabled query boolean false "Camera enabled state"
// @Param created_from query string false "Inclusive creation time lower bound"
// @Param created_to query string false "Inclusive creation time upper bound"
// @Success 200 {integer} integer "Number of cameras"
// @Failure 400 "Invalid filter parameters"
// @Failure 500 "Internal server error"
// @Router /api/camera/count [get]
// @Security BearerAuth
// @Tags Camera
func (c *CameraController) GetCameraCount(ctx *gin.Context) {
	filter, err := helpers.ExtractFilterFromQueryParams(ctx.Request.URL.Query())
	if err != nil {
		log.Printf("Invalid filter: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	count, err := c.cameraService.GetCameraCount(filter)
	if err != nil {
		log.Printf("Error get camera count: %v", err)
		writeCameraError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, count)
}

// UpdateCamera godoc
// @Summary Update camera by ID.
// @Description Updates information for an existing camera record
// @Accept  json
// @Produce json
// @Param id path integer true "Camera ID"
// @Param request body entities.CameraUpdate true "camera entity"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500 "Internal server error"
// @Router /api/camera/{id} [put]
// @Security BearerAuth
// @Tags Camera
func (c *CameraController) UpdateCamera(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Error parse id: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	camera := &entities.CameraUpdate{}
	err = ctx.ShouldBindJSON(camera)
	if err != nil || (camera.Comment == nil && camera.IP == nil && camera.Port == nil && camera.StreamURL == nil && camera.IsEnabled == nil) ||
		(camera.IP != nil && strings.TrimSpace(*camera.IP) == "") {
		log.Printf("Error shouldBindJSON: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = c.cameraService.UpdateCamera(id, camera)
	if err != nil {
		log.Printf("Error update camera: %v", err)
		writeCameraError(ctx, err)
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
// @Failure 409 "Camera has events"
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
		writeCameraError(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func writeCameraError(ctx *gin.Context, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		ctx.Status(http.StatusNotFound)
		return
	}
	var pgError *pq.Error
	if errors.As(err, &pgError) && pgError.Code == "23503" {
		ctx.JSON(http.StatusConflict, gin.H{"error": "camera_has_events"})
		return
	}
	ctx.Status(http.StatusInternalServerError)
}
