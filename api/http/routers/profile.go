package routers

import (
	"build-service-gin/api/http/handlers"
	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	router    *gin.Engine
	clientSys *gin.RouterGroup
	handlers  *handlers.ProfileHandler
}

func NewProfileController(router *gin.Engine, handlers *handlers.ProfileHandler) *ProfileController {
	return &ProfileController{
		router:    router,
		clientSys: router.Group(prefixSystemPath),
		handlers:  handlers,
	}
}

func (app *ProfileController) SetupProfileRoutes() {
	// Set up router for profile
	app.SetupRouterProfile()
}

func (app *ProfileController) SetupRouterProfile() {
	profile := app.clientSys.Group(prefixProfile)
	profile.GET(prefixUserTransactionHistoryPath, app.handlers.GetUserTransactionHistory)
	profile.POST(prefixUserTransactionHistoryPath, app.handlers.CreateUserTransactionHistory)
	profile.PUT(prefixUserTransactionHistoryPath, app.handlers.UpdateUserTransactionHistory)
	profile.DELETE(prefixUserTransactionHistoryPath, app.handlers.DeleteUserTransactionHistory)

	profile.GET(prefixUserTransactionHistoryPostgresPath, app.handlers.GetUserTransactionHistoryPostgres)
	profile.POST(prefixUserTransactionHistoryPostgresPath, app.handlers.CreateUserTransactionHistoryPostgres)
}
