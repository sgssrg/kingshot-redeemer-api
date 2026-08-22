package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"golang.org/x/time/rate"

	echootel "github.com/labstack/echo-opentelemetry"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
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
	// 1. Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// 2. Setup OpenTelemetry Span Exporter
	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		slog.Error("Failed to create OTel Span Exporter", "error", err)
		return
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(spanExporter),
	)
	// IMPORTANT: Register global Tracer Provider & ensure graceful flush on exit
	otel.SetTracerProvider(tp)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error shutting down TracerProvider", "error", err)
		}
	}()

	// 3. Setup OpenTelemetry Metric Exporter
	mReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		slog.Warn("Error creating Metric Reader", "error", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(mReader),
	)
	otel.SetMeterProvider(mp)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mp.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error shutting down MeterProvider", "error", err)
		}
	}()

	e := echo.New()

	// 4. Attach Echo OTel Middleware
	e.Use(echootel.NewMiddlewareWithConfig(echootel.Config{
		ServerName:     "kingshot-api",
		TracerProvider: tp,
		MeterProvider:  mp,

		// Use echov4.Context for echootel callbacks
		Skipper: func(c *echo.Context) bool {
			path := c.Path()
			return path == "/health" || path == "/metrics" || path == "/api/openapi.json"
		},

		OnNextError: func(c *echo.Context, err error) {
			slog.Warn("Handler error", "path", c.Path(), "error", err)
		},

		OnExtractionError: func(c *echo.Context, err error) {
			slog.Warn("Trace context extraction failed", "error", err)
		},

		MetricAttributes: echootel.MetricAttributesFunc(
			func(c *echo.Context, v *echootel.Values) []attribute.KeyValue {
				respCodeStr := strconv.Itoa(v.HTTPResponseStatusCode)
				return []attribute.KeyValue{
					attribute.String("method", v.HTTPMethod),
					attribute.String("route", v.HTTPRoute),
					attribute.String("url_path", v.URLPath),
					attribute.String("client_addr", v.ClientAddress),
					attribute.String("server", v.ServerAddress),
					attribute.String("response_code", respCodeStr),
				}
			},
		),

		SpanStartAttributes: echootel.AttributesFunc(func(c *echo.Context, v *echootel.Values, attr []attribute.KeyValue) []attribute.KeyValue {
			return []attribute.KeyValue{
				attribute.String("service.name", "kingshot-api"),
				attribute.String("deployment.env", "dev"),
			}
		}),
	}))
	// DB & App Routes Setup
	dbApp := db.InitDB()
	defer dbApp.Pool.Close()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/health", online_route.Hello)

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

	e.GET("/api/openapi.json", func(c *echo.Context) error {
		return c.File("docs/swagger.json")
	})

	// 5. Server Startup
	if err := e.Start(":8081"); err != nil && err != http.ErrServerClosed {
		slog.Error("failed to start server", "error", err)
	}
}
