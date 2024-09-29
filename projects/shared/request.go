package shared

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"math/rand"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace" // tracing
)

type RequestOptions struct {
	Method string
	Body   interface{}
}

func Request(ctx context.Context, url string, options *RequestOptions) (*http.Response, error) {
	// simulate a bit more processing time than these toy services perform
	time.Sleep((time.Duration(rand.Intn(50)) * time.Millisecond))
	var (
		req *http.Request
		err error
	)
	if options != nil && options.Body != nil {
		body, err := json.Marshal(options.Body)
		if err != nil {
			return nil, err
		}
		req, _ = http.NewRequestWithContext(ctx, options.Method, url, bytes.NewBuffer(body))
		ctx, req = otelhttptrace.W3C(ctx, req) // tracing
		otelhttptrace.Inject(ctx, req)         // tracing
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
		ctx, req = otelhttptrace.W3C(ctx, req) // tracing
		otelhttptrace.Inject(ctx, req)         // tracing
	}
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(req)
}
