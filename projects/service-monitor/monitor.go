package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"observability/shared"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Target struct {
	ID       string `json:"id"`
	Alerting bool   `json:"alerting"`
}

func newTargets(ctx context.Context) ([]Target, error) {
	endpoint := fmt.Sprintf("%s/v1", registry.Endpoint("panel"))
	resp, err := shared.Request(ctx, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("fetching targets: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading targets: %v", err)
	}
	var targets []Target
	if err := json.Unmarshal(body, &targets); err != nil {
		return nil, fmt.Errorf("parsing targets: %v", err)
	}
	return targets, err
}

func (t Target) getOwnerEmail(ctx context.Context) (string, error) {
	ctx, span := tracer.Start(ctx, "target.getOwnerEmail") // tracing
	defer span.End()                                       // tracing
	endpoint := fmt.Sprintf("%s/v1/%s/owner", registry.Endpoint("panel"), t.ID)
	resp, err := shared.Request(ctx, endpoint, nil)
	if err != nil {
		err := fmt.Errorf("getting owner: %w", err)
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return "", err
	}
	if resp.StatusCode != 200 {
		err := fmt.Errorf("%d response from panel service", resp.StatusCode)
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("reading owner: %s", err)
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return "", err
	}
	return gjson.GetBytes(bytes, "email").String(), nil
}

func (t Target) notify(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "target.notify") // tracing
	defer span.End()                                // tracing
	email, err := t.getOwnerEmail(ctx)
	if err != nil {
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return err
	}
	resp, err := shared.Request(ctx, registry.Endpoint("email")+"/v1", &shared.RequestOptions{
		Method: "POST",
		Body: map[string]string{
			"to":      email,
			"from":    "alarmsystem@test.com",
			"subject": "alerming",
			"body":    "wat",
		},
	})
	if err != nil {
		err := fmt.Errorf("sending email: %w", err)
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return err
	}
	if resp.StatusCode != 200 {
		err := fmt.Errorf("%d response from email service", resp.StatusCode)
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return err
	}
	return nil
}

func checkAllSystems(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "checkAll") // tracing
	defer span.End()                           // tracing
	targets, err := newTargets(ctx)
	if err != nil {
		span.RecordError(err, trace.WithStackTrace(true)) // tracing
		span.SetStatus(codes.Error, err.Error())          // tracing
		return err
	}
	for _, target := range targets {
		if target.Alerting {
			if err := target.notify(ctx); err != nil {
				span.RecordError(err, trace.WithStackTrace(true)) // tracing
				span.SetStatus(codes.Error, err.Error())          // tracing
				return err
			}
		}
	}
	return nil
}
