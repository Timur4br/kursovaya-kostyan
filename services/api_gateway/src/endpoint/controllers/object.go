package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"api_gateway/entities"
	"api_gateway/helpers"
	"api_gateway/services"
)

type ObjectController struct {
	objectService *services.ObjectService
}

func NewObjectController(service *services.ObjectService) *ObjectController {
	return &ObjectController{
		objectService: service,
	}
}

// CreateObject godoc
// @Summary Create object with passed properties.
// @Description Creates a new object in the system
// @Param request body entities.Object true "object entity"
// @Success 201 {object} int
// @Failure 400
// @Failure 500
// @Router /api/object [post]
// @Security BearerAuth
// @Tags Object
func (c *ObjectController) CreateObject(ctx *gin.Context) {
	object := &entities.Object{}
	err := ctx.ShouldBindJSON(&object)
	if err != nil {
		log.Printf("Error shouldBindJSON: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	id, err := c.objectService.CreateObject(object)
	if err != nil {
		log.Printf("Error create object: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, id)
}

// GetObject godoc
// @Summary Get object info by ID.
// @Produce json
// @Param id path integer true "Object ID"
// @Success 200 {object} entities.Object
// @Failure 400
// @Failure 404
// @Failure 500 "Internal server error"
// @Router /api/object/{id} [get]
// @Security BearerAuth
// @Tags Object
func (c *ObjectController) GetObjectByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Error parse id: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	object, err := c.objectService.GetObjectByID(id)
	if err != nil {
		log.Printf("Error get object: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, object)
}

// GetObjectList godoc
// @Summary Get paginated list of objects
// @Produce json
// @Param object_id query integer false "Object ID"
// @Param type query string false "Exact object type"
// @Param created_from query string false "Inclusive creation time lower bound"
// @Param created_to query string false "Inclusive creation time upper bound"
// @Description Retrieves a paginated list of objects
// @Param limit query integer false "limit"
// @Param offset query integer false "offset"
// @Success 200 {array} entities.Object "List of objects"
// @Failure 400 "Invalid pagination parameters"
// @Failure 500 "Internal server error"
// @Router /api/object/list [get]
// @Security BearerAuth
// @Tags Object
func (c *ObjectController) GetObjectList(ctx *gin.Context) {
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

	objects, err := c.objectService.GetObjectList(int(limit), int(offset), filter)
	if err != nil {
		log.Printf("Error get object list: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, objects)
}

// GetCount godoc
// @Summary Get objects count
// @Description Returns the total number of objects in the system
// @Produce json
// @Param object_id query integer false "Object ID"
// @Param type query string false "Exact object type"
// @Param created_from query string false "Inclusive creation time lower bound"
// @Param created_to query string false "Inclusive creation time upper bound"
// @Success 200 {integer} integer "Number of objects"
// @Failure 400 "Invalid filter parameters"
// @Failure 500 "Internal server error"
// @Router /api/object/count [get]
// @Security BearerAuth
// @Tags Object
func (c *ObjectController) GetObjectCount(ctx *gin.Context) {
	filter, err := helpers.ExtractFilterFromQueryParams(ctx.Request.URL.Query())
	if err != nil {
		log.Printf("Invalid filter: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	count, err := c.objectService.GetObjectCount(filter)
	if err != nil {
		log.Printf("Error get object count: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, count)
}

// UpdateObject godoc
// @Summary Update object by ID.
// @Description Updates information for an existing object record
// @Accept  json
// @Produce json
// @Param id path integer true "Object ID"
// @Param request body entities.Object true "object entity"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500 "Internal server error"
// @Router /api/object/{id} [put]
// @Security BearerAuth
// @Tags Object
func (c *ObjectController) UpdateObject(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Error parse id: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	object := &entities.Object{}
	err = ctx.ShouldBindJSON(&object)
	if err != nil {
		log.Printf("Error shouldBindJSON: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = c.objectService.UpdateObject(id, object)
	if err != nil {
		log.Printf("Error update object: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteObject godoc
// @Summary Delete object object by ID.
// @Accept  json
// @Produce json
// @Param id path integer true "Object ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500 "Internal server error"
// @Router /api/object/{id} [delete]
// @Security BearerAuth
// @Tags Object
func (c *ObjectController) DeleteObject(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Error parse id: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = c.objectService.DeleteObject(id)
	if err != nil {
		log.Printf("Error delete object: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
