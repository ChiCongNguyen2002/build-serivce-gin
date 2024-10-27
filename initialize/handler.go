package initialize

import (
	"build-service-gin/api/http/handlers"
	"build-service-gin/api/msgbroker/consumer/core_handle_point"
	"build-service-gin/api/msgbroker/consumer/order"
)

type Handlers struct {
	ProfileHandler   *handlers.ProfileHandler
	PointHandler     *handlers.PointHandler
	OrderHandler     *order.OrderHandler
	CorePointHandler *core_handle_point.CorePointHandler
}

func NewHandlers(services *Services) *Handlers {

	profileHandler := handlers.NewProfileHandler(
		services.profileService,
	)

	pointHandler := handlers.NewPointHandler(
		services.pointService)

	orderHandler := order.NewOrderHandler(
		services.profileService,
	)

	corePointHandler := core_handle_point.NewCorePointHandler(
		services.profileService,
	)

	return &Handlers{
		ProfileHandler:   profileHandler,
		PointHandler:     pointHandler,
		OrderHandler:     orderHandler,
		CorePointHandler: corePointHandler,
	}
}
