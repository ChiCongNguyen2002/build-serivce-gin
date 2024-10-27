package routers

import (
	"build-service-gin/api/http/handlers"
	"github.com/gin-gonic/gin"
)

type PointController struct {
	router    *gin.Engine
	clientSys *gin.RouterGroup
	handlers  *handlers.PointHandler
}

func NewPointController(router *gin.Engine, handlers *handlers.PointHandler) *PointController {
	return &PointController{
		router:    router,
		clientSys: router.Group(prefixSystemPath),
		handlers:  handlers,
	}
}

func (app *PointController) SetupPointRoutes() {
	app.SetupRouterPoint()
}

func (app *PointController) SetupRouterPoint() {
	profile := app.clientSys.Group(prefixPoint)
	profile.POST(prefixPointTransactionPath, app.handlers.CreatePointTransaction)
}
