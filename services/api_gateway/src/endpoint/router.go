package endpoint

import (
	"database/sql"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"

	_ "api_gateway/docs"

	pgDAO "api_gateway/daos/postgres"

	"api_gateway/endpoint/controllers"
	"api_gateway/services"
)

type Router struct {
	Router *gin.Engine
}

func SwaggerHandler(c *gin.Context) {
	doc, err := swag.ReadDoc()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, doc)
}

func NewRouter(db *sql.DB) Router {
	// DAOs
	cameraDAO := pgDAO.NewCameraDAO(db)
	objectDAO := pgDAO.NewObjectDAO(db)
	eventDAO := pgDAO.NewEventDAO(db)

	// Services
	cameraService := services.NewCameraService(cameraDAO)
	objectService := services.NewObjectService(objectDAO)
	eventService := services.NewEventService(eventDAO)

	// Controllers
	cameraController := controllers.NewCameraController(cameraService)
	objectController := controllers.NewObjectController(objectService)
	eventController := controllers.NewEventController(eventService)

	// Gin Settings
	gin.DefaultWriter = log.Writer()
	router := gin.New()

	// Middlewares
	router.Use(LoggingMiddleware(), RecoveryToLogs())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Main Groups
	apiGroup := router.Group("/api")

	// Camera Routes
	cameraGroup := apiGroup.Group("/camera")
	{
		cameraGroup.POST("", cameraController.CreateCamera)
		cameraGroup.GET("/:id", cameraController.GetCameraByID)
		cameraGroup.GET("/count", cameraController.GetCameraCount)
		cameraGroup.GET("/list", cameraController.GetCameraList)
		cameraGroup.PUT("/:id", cameraController.UpdateCamera)
		cameraGroup.DELETE("/:id", cameraController.DeleteCamera)
	}

	// Object Routes
	objectGroup := apiGroup.Group("/object")
	{
		objectGroup.POST("", objectController.CreateObject)
		objectGroup.GET("/:id", objectController.GetObjectByID)
		objectGroup.GET("/count", objectController.GetObjectCount)
		objectGroup.GET("/list", objectController.GetObjectList)
		objectGroup.PUT("/:id", objectController.UpdateObject)
		objectGroup.DELETE("/:id", objectController.DeleteObject)
	}

	// Event Routes
	eventGroup := apiGroup.Group("/event")
	{
		eventGroup.POST("", eventController.CreateEvent)
		eventGroup.GET("/:uuid", eventController.GetEventByUUID)
		eventGroup.GET("/count", eventController.GetEventCount)
		eventGroup.GET("/list", eventController.GetEventList)
		eventGroup.PUT("/:uuid", eventController.UpdateEvent)
		eventGroup.DELETE("/:uuid", eventController.DeleteEvent)
	}

	healthGroup := router.Group("")
	{
		healthGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})
	}
	return Router{Router: router}
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("%s %s", c.Request.Method, c.Request.URL.RequestURI())
		c.Next()
		log.Printf("Response status code: %d", c.Writer.Status())
	}
}

func RecoveryToLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic was recovered: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
			}
		}()
		c.Next()
	}
}
