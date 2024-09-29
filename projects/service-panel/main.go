package main

import (
	"fmt"
	"log"
	"net/http"
	"observability/shared"
	"time"

	"math/rand"

	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const serviceName = "panel"

var tracer = otel.Tracer(serviceName) // tracing

var registry = shared.Registry{}

func init() {
	registry = shared.RegistryFromEnv(serviceName, []string{"email", "user"})
}

func main() {
	log.Printf("%s service starting...\n", serviceName)
	service, port := shared.Server(serviceName)
	service.GET("/v1", func(c echo.Context) error {
		return c.JSON(http.StatusOK, database)
	})
	alarm := service.Group("/v1/:id", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Param("id")
			span := trace.SpanFromContext(c.Request().Context())  // tracing
			span.SetAttributes(attribute.String("system.id", id)) // tracing
			if len(id) != 36 {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": fmt.Sprintf("invalid uuid: %s", id),
				})
			}
			system := systemById(c.Request().Context(), id)
			if system == nil {
				return c.JSON(http.StatusNotFound, map[string]string{
					"error": fmt.Sprintf("system not found: %s", id),
				})
			}
			// simulate specific system experiencing degraded service
			if system.ID == "5a0ab2c1-9f25-4e5d-b578-e9f7401e93d1" {
				time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
			}
			c.Set("system", system)
			owner, err := system.getOwner(c.Request().Context())
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
			}
			span.SetAttributes(attribute.String("user.id", gjson.GetBytes(owner, "id").String()))           // tracing
			span.SetAttributes(attribute.String("user.country", gjson.GetBytes(owner, "country").String())) // tracing
			c.Set("owner", owner)
			return next(c)
		}
	})
	alarm.GET("", func(c echo.Context) error {
		system := c.Get("system").(*System)
		return c.JSON(http.StatusOK, system)
	})
	alarm.POST("/arm", func(c echo.Context) error {
		system := c.Get("system").(*System)
		system.Arm(c.Request().Context())
		return c.JSON(http.StatusOK, system)
	})
	alarm.POST("/disarm", func(c echo.Context) error {
		system := c.Get("system").(*System)
		system.Disarm(c.Request().Context())
		return c.JSON(http.StatusOK, system)
	})
	alarm.POST("/trigger", func(c echo.Context) error {
		system := c.Get("system").(*System)
		system.Trigger(c.Request().Context())
		return c.JSON(http.StatusOK, system)
	})
	alarm.POST("/clear", func(c echo.Context) error {
		system := c.Get("system").(*System)
		system.Clear(c.Request().Context())
		return c.JSON(http.StatusOK, system)
	})
	alarm.GET("/owner", func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.Get("owner"))
	})
	log.Printf("%s service listening on port %s...\n", serviceName, port)
	if err := service.Start(":" + port); err != nil {
		log.Fatalf("%s service startup failed: %v", serviceName, err)
	}
}
