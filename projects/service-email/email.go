package main

import (
	"context"
)

type Email struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (e *Email) Send(ctx context.Context) error {
	_, span := tracer.Start(ctx, "Email.Send") // tracing
	defer span.End()                           // tracing
	//time.Sleep((300 * time.Millisecond) + (time.Duration(rand.Intn(300)) * time.Millisecond))
	return nil
}
