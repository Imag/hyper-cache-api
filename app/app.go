package app

import (
	"github.com/Imag/hyper-cache-api/app/cache"
	"github.com/Imag/hyper-cache-api/app/external"
	"github.com/Imag/hyper-cache-api/app/types"
	"github.com/labstack/echo/v4"
)

type App struct {	
	e *echo.Echo
	CacheService *cache.CacheService
	RetrieveLicense func(key string) (types.License, error)
}

func New() *App {
	app := &App{
		e: echo.New(),
		CacheService: cache.NewCacheService(),
		RetrieveLicense: external.RetrieveLicense,
	}

	app.e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return cacheMiddleware(app, next)
	})

	api := app.e.Group("/api")

	app.registerRoutes(api)

	return app
}

func (a *App) registerRoutes(g *echo.Group) {
	g.GET("/auth/:license", func(c echo.Context) error {
		return c.JSON(200, c.Get("license"))
	})
}

func (a *App) Run() {
	go a.CacheService.RunCacheHandler()
	a.e.Start(":3000")
}