package http

import (
	"build-service-gin/api/http/middlewares"
	"build-service-gin/api/http/routers"
	"build-service-gin/pkg/helpers/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

const prefixPath = "build-service-gin"
const healthPath = prefixPath + "/v1/health"

func (app *httpServ) InitRouters(g *gin.Engine) {
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	g.Use(middlewares.AddExtraDataForRequestContext)
	g.Use(middlewares.Logging)

	g.GET(healthPath, func(c *gin.Context) {
		if healthCheck {
			c.JSON(http.StatusOK, resp.BuildSuccessResp(resp.LangEN, nil))
		} else {
			c.JSON(http.StatusInternalServerError, resp.BuildErrorResp(resp.ErrSystem, "", resp.LangEN))
		}
	})

	//profile router
	controller := routers.NewProfileController(g, app.profileHandler)
	controller.SetupProfileRoutes()

	//point router
	pointController := routers.NewPointController(g, app.pointHandler)
	pointController.SetupPointRoutes()
}
