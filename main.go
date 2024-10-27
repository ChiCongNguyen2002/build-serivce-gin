package main

import (
	apiHttp "build-service-gin/api/http"
	"build-service-gin/api/http/middlewares"
	"build-service-gin/common/logger"
	"build-service-gin/common/mongodb"
	postgres "build-service-gin/common/postgresql"
	"build-service-gin/config"
	"build-service-gin/initialize"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.InitLog(os.Getenv("SERVICE_ID"))
	log := logger.GetLogger()
	log.Info().Any("service", os.Getenv("SERVICE_ID")).Msg("Start services")

	// Load configuration
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("load config fail! %s", err)
	}

	// Set health check status
	apiHttp.SetHealthCheck(true)
	g := gin.Default()

	// Initialize metrics
	middlewares.InitMetrics()
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Create context for handling shutdown signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Connect to MongoDB
	dbStorage, err := mongodb.ConnectMongoDB(context.Background(), &conf.MongoDBConfig)
	if err != nil {
		log.Fatal().Msgf("connect mongodb failed! %s", err)
	}

	// Connect to Postgres
	postgresql, err := postgres.ConnectPostgresql(context.Background(), &conf.PostgresConfig)
	if err != nil {
		log.Fatal().Msgf("connect postgresql failed! %s", err)
	}

	//redisClient, err := redis.ConnectRedis(context.Background(), conf.RedisConfig)
	//if err != nil {
	//	log.Fatal().Err(err).Msg("Initialize redis client failed.")
	//}

	// Initialize clients
	clients := initialize.NewClients()

	// Initialize repositories
	repo := initialize.NewRepositories(dbStorage, postgresql)

	// Initialize services
	//service := initialize.NewServices(conf, clients, repo, redisClient)
	service := initialize.NewServices(conf, clients, repo)
	// Initialize handlers
	handler := initialize.NewHandlers(service)

	// Uncomment if you're using a message broker
	// go func() {
	// 	msgBroker := msgbroker.NewMsgBroker(conf, handler.OrderHandler, handler.CorePointHandler)
	// 	msgBroker.Start(ctx)
	// }()

	// Create HTTP server instance
	srv := apiHttp.NewHttpServe(conf, handler.ProfileHandler, handler.PointHandler)
	srv.Start(g)

	// Handle graceful shutdown
	<-ctx.Done()
	apiHttp.SetHealthCheck(false)

	// Set a timeout for shutdown
	cancelCtx, cc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cc()

	// Shutdown the server
	if err = srv.Shutdown(cancelCtx); err != nil {
		log.Fatal().Msgf("force shutdown services: %v", err)
	}

	log.Info().Msg("Server exited gracefully")
}
