package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"api_gateway/entities"
	"api_gateway/helpers"
	"api_gateway/services"
)

type EventController struct {
	eventService *services.EventService
}

func NewEventController(service *services.EventService) *EventController {
	return &EventController{
		eventService: service,
	}
}

// CreateEvent godoc
// @Summary Create event with passed properties.
// @Description Creates a new event in the system
// @Param request body entities.Event true "event entity"
// @Success 201 {string} string
// @Failure 400
// @Failure 500
// @Router /api/event [post]
// @Security BearerAuth
// @Tags Event
func (c *EventController) CreateEvent(ctx *gin.Context) {
	event := &entities.Event{}
	err := ctx.ShouldBindJSON(&event)
	if err != nil {
		log.Printf("Error shouldBindJSON: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	uuid, err := c.eventService.CreateEvent(event)
	if err != nil {
		log.Printf("Error create event: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, uuid)
}

// GetEvent godoc
// @Summary Get event info by UUID.
// @Produce json
// @Param uuid path string true "Event UUID"
// @Success 200 {object} entities.Event
// @Failure 400
// @Failure 404
// @Failure 500 "Internal server error"
// @Router /api/event/{uuid} [get]
// @Security BearerAuth
// @Tags Event
func (c *EventController) GetEventByUUID(ctx *gin.Context) {
	uuid, err := extractEventUUID(ctx)
	if err != nil {
		log.Printf("Error parse uuid: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	event, err := c.eventService.GetEventByUUID(uuid)
	if err != nil {
		log.Printf("Error get event: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, event)
}

// GetEventList godoc
// @Summary Get paginated list of events
// @Produce json
// @Param camera_id query integer false "Camera ID"
// @Param object_id query integer false "Object ID"
// @Param type query string false "Exact event type"
// @Param confidence_min query number false "Inclusive confidence lower bound"
// @Param confidence_max query number false "Inclusive confidence upper bound"
// @Param created_from query string false "Inclusive creation time lower bound"
// @Param created_to query string false "Inclusive creation time upper bound"
// @Description Retrieves a paginated list of events
// @Param limit query integer false "limit"
// @Param offset query integer false "offset"
// @Success 200 {array} entities.Event "List of events"
// @Failure 400 "Invalid pagination parameters"
// @Failure 500 "Internal server error"
// @Router /api/event/list [get]
// @Security BearerAuth
// @Tags Event
func (c *EventController) GetEventList(ctx *gin.Context) {
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

	events, err := c.eventService.GetEventList(int(limit), int(offset), filter)
	if err != nil {
		log.Printf("Error get event list: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, events)
}

// GetCount godoc
// @Summary Get events count
// @Description Returns the total number of events in the system
// @Produce json
// @Param camera_id query integer false "Camera ID"
// @Param object_id query integer false "Object ID"
// @Param type query string false "Exact event type"
// @Param confidence_min query number false "Inclusive confidence lower bound"
// @Param confidence_max query number false "Inclusive confidence upper bound"
// @Param created_from query string false "Inclusive creation time lower bound"
// @Param created_to query string false "Inclusive creation time upper bound"
// @Success 200 {integer} integer "Number of events"
// @Failure 400 "Invalid filter parameters"
// @Failure 500 "Internal server error"
// @Router /api/event/count [get]
// @Security BearerAuth
// @Tags Event
func (c *EventController) GetEventCount(ctx *gin.Context) {
	filter, err := helpers.ExtractFilterFromQueryParams(ctx.Request.URL.Query())
	if err != nil {
		log.Printf("Invalid filter: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	count, err := c.eventService.GetEventCount(filter)
	if err != nil {
		log.Printf("Error get event count: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, count)
}

// UpdateEvent godoc
// @Summary Update event by UUID.
// @Description Updates information for an existing event record
// @Accept  json
// @Produce json
// @Param uuid path string true "Event UUID"
// @Param request body entities.Event true "event entity"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500 "Internal server error"
// @Router /api/event/{uuid} [put]
// @Security BearerAuth
// @Tags Event
func (c *EventController) UpdateEvent(ctx *gin.Context) {
	uuid, err := extractEventUUID(ctx)
	if err != nil {
		log.Printf("Error parse uuid: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	event := &entities.Event{}
	err = ctx.ShouldBindJSON(&event)
	if err != nil {
		log.Printf("Error shouldBindJSON: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = c.eventService.UpdateEvent(uuid, event)
	if err != nil {
		log.Printf("Error update event: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteEvent godoc
// @Summary Delete event by UUID.
// @Accept  json
// @Produce json
// @Param uuid path string true "Event UUID"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500 "Internal server error"
// @Router /api/event/{uuid} [delete]
// @Security BearerAuth
// @Tags Event
func (c *EventController) DeleteEvent(ctx *gin.Context) {
	uuid, err := extractEventUUID(ctx)
	if err != nil {
		log.Printf("Error parse uuid: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = c.eventService.DeleteEvent(uuid)
	if err != nil {
		log.Printf("Error delete event: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func extractEventUUID(ctx *gin.Context) (string, error) {
	var params struct {
		UUID string `uri:"uuid" binding:"required,uuid"`
	}
	err := ctx.ShouldBindUri(&params)
	return params.UUID, err
}
