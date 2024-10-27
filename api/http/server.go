package http

import (
	"build-service-gin/api/http/handlers"
	"build-service-gin/common/logger"
	"build-service-gin/config"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var (
	healthCheck bool
	mu          sync.RWMutex
)

func SetHealthCheck(status bool) {
	mu.Lock()
	defer mu.Unlock()
	healthCheck = status
}

type HttpServInterface interface {
	Start(g *gin.Engine)
	Shutdown(ctx context.Context) error
}

type httpServ struct {
	conf           *config.SystemConfig
	profileHandler *handlers.ProfileHandler
	pointHandler   *handlers.PointHandler
	httpServer     *http.Server
	//coreHandler    *order.OrderHandler
	//earnHandler    *core_handle_point.CorePointHandler
}

func NewHttpServe(
	conf *config.SystemConfig,
	profileHandler *handlers.ProfileHandler,
	pointHandler *handlers.PointHandler,
	// coreHandler *order.OrderHandler,
	// earnHandler *core_handle_point.CorePointHandler,
) *httpServ {
	return &httpServ{
		conf:           conf,
		profileHandler: profileHandler,
		pointHandler:   pointHandler,
		httpServer: &http.Server{ // Initialize the HTTP server
			Addr: fmt.Sprintf(":%d", conf.HttpPort),
		},
		//coreHandler:    coreHandler,
		//earnHandler:    earnHandler,
	}
}

func (app *httpServ) Start(g *gin.Engine) {
	log := logger.GetLogger()
	app.InitRouters(g)
	httpPort := app.conf.HttpPort
	go func() {
		err := g.Run(fmt.Sprintf(":%d", httpPort))
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("can't start gin: %v", err)
		}
	}()
	log.Info().Msg("HTTP server started on port: " + fmt.Sprintf("%d", httpPort))
}

func (app *httpServ) Shutdown(ctx context.Context) error {
	log := logger.GetLogger()
	if err := app.httpServer.Shutdown(ctx); err != nil {
		log.Error().Msgf("Server shutdown failed: %v", err)
		return err
	}
	log.Info().Msg("Server shutdown gracefully")
	return nil
}
