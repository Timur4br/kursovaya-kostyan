package httpapi

import (
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/handlers"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Controllers struct {
	User            handlers.UserController
	Location        handlers.LocationController
	Camera          handlers.CameraController
	Person          handlers.PersonController
	ObjectType      handlers.ObjectTypeController
	Object          handlers.ObjectController
	EventType       handlers.EventTypeController
	Event           handlers.EventController
	VideoSegment    handlers.VideoSegmentController
	Trajectory      handlers.TrajectoryController
	TrajectoryPoint handlers.TrajectoryPointController
	QueryHistory    handlers.QueryHistoryController
	QueryEvent      handlers.QueryEventController
}

func registerEntities(router *gin.Engine, dependencies *Controllers, logger *slog.Logger) {
	userHandler := handlers.NewUserHandler(dependencies.User, logger)
	router.POST("/user", userHandler.Create)
	router.GET("/user/:id", userHandler.Get)
	router.PUT("/user/:id", userHandler.Update)
	router.DELETE("/user/:id", userHandler.Delete)
	locationHandler := handlers.NewLocationHandler(dependencies.Location, logger)
	router.POST("/location", locationHandler.Create)
	router.GET("/location/:id", locationHandler.Get)
	router.PUT("/location/:id", locationHandler.Update)
	router.DELETE("/location/:id", locationHandler.Delete)
	cameraHandler := handlers.NewCameraHandler(dependencies.Camera, logger)
	router.POST("/camera", cameraHandler.Create)
	router.GET("/camera/:id", cameraHandler.Get)
	router.PUT("/camera/:id", cameraHandler.Update)
	router.DELETE("/camera/:id", cameraHandler.Delete)
	personHandler := handlers.NewPersonHandler(dependencies.Person, logger)
	router.POST("/person", personHandler.Create)
	router.GET("/person/:id", personHandler.Get)
	router.PUT("/person/:id", personHandler.Update)
	router.DELETE("/person/:id", personHandler.Delete)
	object_typeHandler := handlers.NewObjectTypeHandler(dependencies.ObjectType, logger)
	router.POST("/object_type", object_typeHandler.Create)
	router.GET("/object_type/:id", object_typeHandler.Get)
	router.PUT("/object_type/:id", object_typeHandler.Update)
	router.DELETE("/object_type/:id", object_typeHandler.Delete)
	objectHandler := handlers.NewObjectHandler(dependencies.Object, logger)
	router.POST("/object", objectHandler.Create)
	router.GET("/object/:id", objectHandler.Get)
	router.PUT("/object/:id", objectHandler.Update)
	router.DELETE("/object/:id", objectHandler.Delete)
	event_typeHandler := handlers.NewEventTypeHandler(dependencies.EventType, logger)
	router.POST("/event_type", event_typeHandler.Create)
	router.GET("/event_type/:id", event_typeHandler.Get)
	router.PUT("/event_type/:id", event_typeHandler.Update)
	router.DELETE("/event_type/:id", event_typeHandler.Delete)
	eventHandler := handlers.NewEventHandler(dependencies.Event, logger)
	router.POST("/event", eventHandler.Create)
	router.GET("/event/:id", eventHandler.Get)
	router.PUT("/event/:id", eventHandler.Update)
	router.DELETE("/event/:id", eventHandler.Delete)
	video_segmentHandler := handlers.NewVideoSegmentHandler(dependencies.VideoSegment, logger)
	router.POST("/video_segment", video_segmentHandler.Create)
	router.GET("/video_segment/:id", video_segmentHandler.Get)
	router.PUT("/video_segment/:id", video_segmentHandler.Update)
	router.DELETE("/video_segment/:id", video_segmentHandler.Delete)
	trajectoryHandler := handlers.NewTrajectoryHandler(dependencies.Trajectory, logger)
	router.POST("/trajectory", trajectoryHandler.Create)
	router.GET("/trajectory/:id", trajectoryHandler.Get)
	router.PUT("/trajectory/:id", trajectoryHandler.Update)
	router.DELETE("/trajectory/:id", trajectoryHandler.Delete)
	trajectory_pointHandler := handlers.NewTrajectoryPointHandler(dependencies.TrajectoryPoint, logger)
	router.POST("/trajectory_point", trajectory_pointHandler.Create)
	router.GET("/trajectory_point/:id", trajectory_pointHandler.Get)
	router.PUT("/trajectory_point/:id", trajectory_pointHandler.Update)
	router.DELETE("/trajectory_point/:id", trajectory_pointHandler.Delete)
	query_historyHandler := handlers.NewQueryHistoryHandler(dependencies.QueryHistory, logger)
	router.POST("/query_history", query_historyHandler.Create)
	router.GET("/query_history/:id", query_historyHandler.Get)
	router.PUT("/query_history/:id", query_historyHandler.Update)
	router.DELETE("/query_history/:id", query_historyHandler.Delete)
	query_eventHandler := handlers.NewQueryEventHandler(dependencies.QueryEvent, logger)
	router.POST("/query_event", query_eventHandler.Create)
	router.GET("/query_event/:query_id/:event_id", query_eventHandler.Get)
	router.PUT("/query_event/:query_id/:event_id", query_eventHandler.Update)
	router.DELETE("/query_event/:query_id/:event_id", query_eventHandler.Delete)
}
