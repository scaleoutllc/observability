package main

import (
	"log"
	"net/http"
	"observability/shared"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

var serviceName = "monitor"

var tracer = otel.Tracer(serviceName) // tracing

var registry = shared.Registry{}

func init() {
	registry = shared.RegistryFromEnv(serviceName, []string{"email", "panel"})
}

func main() {
	log.Printf("%s service starting...\n", serviceName)
	service, port := shared.Server(serviceName)
	service.GET("/v1/all", func(c echo.Context) error {
		checkAllSystems(c.Request().Context())
		return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
	log.Printf("%s service listening on port %s...\n", serviceName, port)
	if err := service.Start(":" + port); err != nil {
		log.Fatalf("%s service startup failed: %v", serviceName, err)
	}
}
