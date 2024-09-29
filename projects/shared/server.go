package shared

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho" // tracing
	"go.opentelemetry.io/otel"                                                      // tracing
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"                             // tracing
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"               // tracing
	"go.opentelemetry.io/otel/propagation"                                          // tracing
	sdktrace "go.opentelemetry.io/otel/sdk/trace"                                   // tracing
)

func Server(name string) (*echo.Echo, string) {
	field := fmt.Sprintf("%s_SERVER_PORT", strings.ToUpper(name))
	port := os.Getenv(field)
	if port == "" {
		log.Fatalf("env %s must be specified", field)
	}
	e := echo.New()
	StartTelemetry()                 // tracing
	e.Use(otelecho.Middleware(name)) // tracing
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/liveness", func(c echo.Context) error {
		return c.String(http.StatusOK, http.StatusText(http.StatusOK))
	})
	e.GET("/readiness", func(c echo.Context) error {
		return c.String(http.StatusOK, http.StatusText(http.StatusOK))
	})
	return e, port
}

func StartTelemetry() { // tracing
	client := otlptracehttp.NewClient()                     // tracing
	exp, err := otlptrace.New(context.Background(), client) // tracing
	if err != nil {                                         // tracing
		log.Fatalf("failed to initialize telemetry exporter: %e", err) // tracing
	} // tracing
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp)) // tracing
	otel.SetTracerProvider(tp)                                  // tracing
	otel.SetTextMapPropagator(                                  // tracing
		propagation.NewCompositeTextMapPropagator( // tracing
			propagation.TraceContext{}, // tracing
			propagation.Baggage{},      // tracing
		), // tracing
	) // tracing
} // tracing
