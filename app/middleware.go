package app

import (
	"net/http"
	"strings"

	"github.com/Imag/hyper-cache-api/app/cache"
	"github.com/Imag/hyper-cache-api/app/types"
	"github.com/labstack/echo/v4"
)

func cacheMiddleware(app *App, next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !strings.HasPrefix(c.Path(), "/api/auth/:license") {
			return next(c)
		}

		var license types.License
		cacheExists, err := app.CacheService.FindCache(cache.LicenseCache, c.Param("license"), &license)
		if err != nil {
			return next(c)
		}

		if cacheExists {
			c.Set("license", license)
		} else {
			l, err := app.RetrieveLicense(c.Param("license"))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err).SetInternal(err)
			}

			if l.Key == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized License")
			}

			c.Set("license", l)

			app.CacheService.UpsertCache(cache.LicenseCache, c.Param("license"), &l)
		}

		return next(c)
	}
}