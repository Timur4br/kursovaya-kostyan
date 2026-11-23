// Package app is the composition root for HTTP, controllers, services and dbDAO.
package app

import (
	"database/sql"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/controllers"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dao"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/httpapi"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/services"
	"log/slog"
	"net/http"
)

func Handler(db *sql.DB, check httpapi.ReadinessCheck, logger *slog.Logger) http.Handler {
	dbDAO := dao.New(db)
	dependencies := &httpapi.Controllers{
		User:            controllers.NewUserController(services.NewUserService(dbDAO)),
		Location:        controllers.NewLocationController(services.NewLocationService(dbDAO)),
		Camera:          controllers.NewCameraController(services.NewCameraService(dbDAO)),
		Person:          controllers.NewPersonController(services.NewPersonService(dbDAO)),
		ObjectType:      controllers.NewObjectTypeController(services.NewObjectTypeService(dbDAO)),
		Object:          controllers.NewObjectController(services.NewObjectService(dbDAO)),
		EventType:       controllers.NewEventTypeController(services.NewEventTypeService(dbDAO)),
		Event:           controllers.NewEventController(services.NewEventService(dbDAO)),
		VideoSegment:    controllers.NewVideoSegmentController(services.NewVideoSegmentService(dbDAO)),
		Trajectory:      controllers.NewTrajectoryController(services.NewTrajectoryService(dbDAO)),
		TrajectoryPoint: controllers.NewTrajectoryPointController(services.NewTrajectoryPointService(dbDAO)),
		QueryHistory:    controllers.NewQueryHistoryController(services.NewQueryHistoryService(dbDAO)),
		QueryEvent:      controllers.NewQueryEventController(services.NewQueryEventService(dbDAO)),
	}
	return httpapi.Handler(check, logger, dependencies)
}
