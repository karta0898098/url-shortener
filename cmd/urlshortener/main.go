package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"url-shortener/cmd/urlshortener/configs"
	"url-shortener/pkg/app/urlshortener/endpoints"
	"url-shortener/pkg/app/urlshortener/repository"
	"url-shortener/pkg/app/urlshortener/service"
	th "url-shortener/pkg/app/urlshortener/transports/http"
	"url-shortener/pkg/bloom"
	"url-shortener/pkg/db"
	"url-shortener/pkg/http"
	"url-shortener/pkg/http/middleware"
	"url-shortener/pkg/logging"
	"url-shortener/pkg/redis"
)

// Application define application
type Application struct {
	logger     zerolog.Logger
	config     configs.Configurations
	httpServer *echo.Echo
	handler    *th.Handler

	Close func()
}

// NewApplication new application
func NewApplication(
	logger zerolog.Logger,
	config configs.Configurations,
) *Application {

	// init database
	dbConn, err := db.NewConnection(config.Database)
	if err != nil {
		logger.
			Panic().
			Err(err).
			Msg("failed to connection database")
	}

	rds, err := redis.NewRedis(config.Redis)
	if err != nil {
		logger.
			Panic().
			Err(err).
			Msg("failed to connection redis")
	}

	bf := bloom.NewRedisFilter(config.BloomFilterNamespace, rds)
	repo := repository.New(dbConn)
	svc := service.New(repo, bf)
	e := endpoints.New(svc)
	h := th.NewHandler(e)

	return &Application{
		logger:     logger,
		config:     config,
		httpServer: http.NewEcho(config.HTTP),
		handler:    h,
		Close: func() {
			rds.Close()
		},
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// read all config
	config := configs.NewConfig("")

	// setup default logger
	logger := logging.Setup(config.Log)

	// make app
	app := NewApplication(logger, config)
	app.httpServer.Pre(middleware.NewLoggerMiddleware(logger))
	app.httpServer.Use(middleware.NewLoggingMiddleware())
	app.httpServer.Use(middleware.RecordErrorMiddleware())
	app.MakeRouter()

	wg := &sync.WaitGroup{}

	go app.startHttpServer(ctx, wg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	cancel()
	wg.Wait()
}

// startHttpServer wrap http start
func (app *Application) startHttpServer(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	app.logger.Info().Msgf("start http server on %v", app.config.HTTP.Port)

	go func() {
		err := app.httpServer.Start(app.config.HTTP.Port)
		if err != nil {
			app.logger.Error().Err(err).Msg("http server shutdown ...")
		}
	}()

	<-ctx.Done()

	// gracefulShutdown wrap graceful shutdown http server
	// timeout 5 sec will direct close server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.httpServer.Shutdown(ctx)
	if err != nil {
		app.logger.
			Error().
			Err(err).
			Msg("failed to http shutdown")
	}

	app.Close()
}

// MakeRouter register router into echo http server
func (app *Application) MakeRouter() *Application {
	app.httpServer.POST("/api/v1/urls", app.handler.ShortURL)
	app.httpServer.GET("/:url", app.handler.RedirectURL)
	return app
}
