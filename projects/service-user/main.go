package main

import (
	"fmt"
	"log"
	"net/http"
	"observability/shared"
	"time"

	"math/rand"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const serviceName = "user"

var tracer = otel.Tracer(serviceName) // tracing

func main() {
	startTime := time.Now()
	log.Printf("%s service starting...\n", serviceName)
	service, port := shared.Server(serviceName)
	service.GET("/v1", func(c echo.Context) error {
		return c.JSON(http.StatusOK, getAllUsers(c.Request().Context()))
	})
	service.GET("/v1/:id", func(c echo.Context) error {
		id := c.Param("id")
		span := trace.SpanFromContext(c.Request().Context()) // tracing
		span.SetAttributes(attribute.String("user.id", id))  // tracing
		if len(id) != 36 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("invalid uuid: %s", id),
			})
		}
		user := getUserByID(c.Request().Context(), id)
		if user == nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": fmt.Sprintf("user not found: %s", id),
			})
		}
		span.SetAttributes(attribute.String("user.country", user.Country)) // tracing
		// simulate country experiencing degraded service
		if user.Country == "Ireland" {
			// for every second of uptime add 10ms of additional delay
			time.Sleep(time.Duration(time.Since(startTime).Seconds()) * (time.Millisecond * 10))
		}
		// simulate specific user experiencing degraded service
		if user.ID == "57efbff5-c520-4d5a-ae30-53b73c6d7f2a" {
			time.Sleep((1 * time.Second) + (time.Duration(rand.Intn(2000)) * time.Millisecond))
		}
		return c.JSON(http.StatusOK, user)
	})
	service.GET("/v1/by-system/:id", func(c echo.Context) error {
		id := c.Param("id")
		span := trace.SpanFromContext(c.Request().Context())  // tracing
		span.SetAttributes(attribute.String("system.id", id)) // tracing
		if len(id) != 36 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("invalid uuid: %s", id),
			})
		}
		user := getUserBySystemID(c.Request().Context(), id)
		if user == nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": fmt.Sprintf("user not found: %s", id),
			})
		}
		span.SetAttributes(attribute.String("user.id", user.ID))           // tracing
		span.SetAttributes(attribute.String("user.country", user.Country)) // tracing
		return c.JSON(http.StatusOK, user)
	})
	log.Printf("%s service listening on port %s...\n", serviceName, port)
	if err := service.Start(":" + port); err != nil {
		log.Fatalf("%s service startup failed: %v", serviceName, err)
	}
}
