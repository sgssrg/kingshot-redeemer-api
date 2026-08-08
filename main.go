// @title Kingshot Redeem API
// @version 1.0
// @description API for redeeming KingShot gift codes on behalf of registered players. OpenAPI spec served at /api/openapi.json for import into GitBook.
// @host localhost:8081
// @BasePath /
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

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
	e := echo.New()
	ctx := context.Background()
	// e.Use(RateLimitMiddleware(0.5, 5))
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

		//Skipper:             nil,
		//OnNextError:         nil,
		//OnExtractionError:   nil,
		//MeterProvider:       nil,
		//Propagators:         nil,
		//SpanStartOptions:    nil,
		//SpanStartAttributes: nil,
		//SpanEndAttributes:   nil,
		//MetricAttributes:    nil,
		//Metrics:             nil,
	}))
	// Initialize DB from database.go
	dbApp := db.InitDB()
	defer dbApp.Pool.Close()

	// Execute your queries
	players, err := dbApp.Queries.GetAllPlayers(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch players: %v", err)
	}

	fmt.Printf("Fetched %d players!\n", len(players))

	// Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", online_route.Hello)

	// Player Routes
	playerRouter := e.Group("/player")
	playerRouter.POST("/add", player.AddPlayer)
	playerRouter.DELETE("/del/:pid", player.DeletePlayer)

	redeemRouter := e.Group("/redeem")
	redeemRouter.POST("/db", redeem.RedeemWithDB, RateLimitMiddleware(0.75, 5))
	redeemRouter.GET("/ws", redeem.WSTest)

	scraperRouter := e.Group("/scraper")
	scraperRouter.GET("/player/:fid", scraper.StratForgePlayerScraperAPI)

	// OpenAPI spec (generated via `swag init`) for GitBook import
	e.GET("/api/openapi.json", func(c *echo.Context) error {
		return c.File("docs/swagger.json")
	})

	// Start server
	if err := e.Start(":8081"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
