package main

import (
	"ecampus/config"
	"ecampus/database"
	"ecampus/internal/services"
	"ecampus/routes"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type App struct {
	DB     *database.DB
	Logger *zap.Logger
	Router *gin.Engine
	Config *config.Config
}

func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	// Initialize user routes
	routes.InitializeUserRoutes(r)
	return r
}

func (app *App) EchoHandler(c *gin.Context) {
	if _, err := io.Copy(c.Writer, c.Request.Body); err != nil {
		app.Logger.Warn("Failed to handle request", zap.Error(err))
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	c.Status(http.StatusOK)
}

func NewApp(logger *zap.Logger, db *database.DB, cfg *config.Config) *App {
	return &App{
		DB:     db,
		Logger: logger,
		Config: cfg,
	}
}

func main() {
	fx.New(
		fx.Provide(
			config.New,
			NewLogger,
			database.New,
			services.NewUserService,
			NewApp,
			NewRouter,
		),
		fx.Invoke(func(app *App, router *gin.Engine) {
			app.Router = router

			// Create database schema
			if err := app.DB.CreateSchema(); err != nil {
				app.Logger.Fatal("Failed to create database schema", zap.Error(err))
			}

			// Start the server
			app.Logger.Info("Starting server", zap.String("addr", app.Config.Port))
			if err := router.Run(app.Config.Port); err != nil {
				app.Logger.Fatal("Failed to start server", zap.Error(err))
			}
		}),
	).Run()
}
