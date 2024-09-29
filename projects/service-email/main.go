package main

import (
	"fmt"
	"log"
	"net/http"

	"observability/shared"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

const serviceName = "email"

var tracer = otel.Tracer(serviceName) // tracing

func main() {
	log.Printf("%s service starting...\n", serviceName)
	service, port := shared.Server(serviceName)
	service.POST("/v1", func(c echo.Context) error {
		email := &Email{}
		if err := c.Bind(email); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"status": "bad request",
			})
		}
		if err := email.Send(c.Request().Context()); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"status": fmt.Sprintf("failed to send email: %s", err),
			})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status": "email sent",
		})
	})
	log.Printf("%s service listening on port %s...\n", serviceName, port)
	if err := service.Start(":" + port); err != nil {
		log.Fatalf("%s service startup failed: %v", serviceName, err)
	}
}
