// @title Kingshot Redeem API
// @version 1.0
// @description API for redeeming KingShot gift codes on behalf of registered players. OpenAPI spec served at /api/openapi.json for import into GitBook.
// @host localhost:8081
// @BasePath /
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"golang.org/x/time/rate"

	echootel "github.com/labstack/echo-opentelemetry"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"gitlab.com/ribonin/apis/kingshot-redeem/db"

	online_route "gitlab.com/ribonin/apis/kingshot-redeem/routes"
	player "gitlab.com/ribonin/apis/kingshot-redeem/routes/Player"
	redeem "gitlab.com/ribonin/apis/kingshot-redeem/routes/Redeem"
	scraper "gitlab.com/ribonin/apis/kingshot-redeem/routes/Scraper"
)

func RateLimitMiddleware(limit rate.Limit, burst int) echo.MiddlewareFunc {
	limiter := rate.NewLimiter(limit, burst)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "429 Too Many Requests",
				})
			}
			return next(c)
		}
	}
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env")
	}
	e := echo.New()
	ctx := context.Background()

	// Create trace exporter using environment variables
	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		// handle error
	}

	// Create trace provider with the exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(spanExporter),
	)

	e.Use(echootel.NewMiddlewareWithConfig(echootel.Config{
		ServerName:     "my-server",
		TracerProvider: tp,
	}))

	// Initialize DB once at startup
	dbApp := db.InitDB()
	defer dbApp.Pool.Close()

	// Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Routes - pass dbApp to handlers
	e.GET("/", online_route.Hello)

	// Player Routes
	playerRouter := e.Group("/player")
	playerRouter.POST("/add", player.AddPlayer(dbApp))
	playerRouter.DELETE("/del/:pid", player.DeletePlayer(dbApp))
	playerRouter.PATCH("/update/all", player.UpdateAllPlayer(dbApp))
	playerRouter.PATCH("/update/:pid", player.UpdatePlayer(dbApp))

	redeemRouter := e.Group("/redeem")
	redeemRouter.POST("/db", redeem.RedeemWithDB(dbApp), RateLimitMiddleware(0.75, 5))
	redeemRouter.POST("/single", redeem.RedeemSingle(dbApp))
	redeemRouter.GET("/ws", redeem.WSTest)

	scraperRouter := e.Group("/scraper")
	scraperRouter.GET("/player/:fid", scraper.StratForgePlayerScraperAPI(dbApp))

	// OpenAPI spec (generated via `swag init`) for GitBook import
	e.GET("/api/openapi.json", func(c *echo.Context) error {
		return c.File("docs/swagger.json")
	})

	// Start server
	if err := e.Start(":8081"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
