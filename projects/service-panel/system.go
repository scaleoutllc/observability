package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"observability/shared"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type System struct {
	ID       string `json:"id"`
	Armed    bool   `json:"armed"`
	Alerting bool   `json:"alerting"`
}

func (s *System) Arm(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "system.Arm") // tracing
	defer span.End()                             // tracing
	if s.Armed {
		err := errors.New("system armed when already armed")
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return err
	}
	s.Armed = true
	return s.emailStateChange(ctx)
}

func (s *System) Disarm(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "system.Disarm") // tracing
	defer span.End()                                // tracing
	if !s.Armed {
		err := errors.New("system disarmed when already disarmed")
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return err
	}
	s.Armed = false
	return s.emailStateChange(ctx)
}

func (s *System) Trigger(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "system.Trigger") // tracing
	defer span.End()                                 // tracing
	s.Alerting = true
	return s.emailStateChange(ctx)
}

func (s *System) Clear(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "system.Clear") // tracing
	defer span.End()                               // tracing
	s.Alerting = false
	return s.emailStateChange(ctx)
}

func (s *System) getOwner(ctx context.Context) ([]byte, error) {
	ctx, span := tracer.Start(ctx, "system.getOwner") // tracing
	defer span.End()                                  // tracing
	endpoint := fmt.Sprintf("%s/v1/by-system/%s", registry.Endpoint("user"), s.ID)
	resp, err := shared.Request(ctx, endpoint, nil)
	if err != nil {
		err := fmt.Errorf("getting alarm system: %w", err)
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err := fmt.Errorf("%d response from user service", resp.StatusCode)
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func (s *System) emailStateChange(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "system.emailStateChange") // tracing
	defer span.End()                                          // tracing
	owner, err := s.getOwner(ctx)
	if err != nil {
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return err
	}
	send, err := shared.Request(ctx, registry.Endpoint("email")+"/v1", &shared.RequestOptions{
		Method: "POST",
		Body: map[string]string{
			"to":      gjson.GetBytes(owner, "email").String(),
			"from":    "system@test.com",
			"subject": "stateChange",
			"body":    "yup",
		},
	})
	if err != nil || send.StatusCode != 200 {
		err := fmt.Errorf("email failure: %s", err)
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return err
	}
	defer send.Body.Close()
	return nil
}
