package tracer

import (
	"context"
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/zipkin"
	_ "github.com/Azure/go-autorest/tracing/opencensus"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
)

func init() {
	localEndpoint, err := openzipkin.NewEndpoint("terraform", "192.168.1.5:5454")
	if err != nil {
		log.Fatalf("Failed to create the local zipkinEndpoint: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	ze := zipkin.NewExporter(reporter, localEndpoint)
	trace.RegisterExporter(ze)

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}

type TraceSpanKey struct{}

type tracer struct{}

func (t *tracer) NewTransport(base *http.Transport) http.RoundTripper {
	return base
}

func (t *tracer) StartSpan(ctx context.Context, name string) context.Context {
	newctx, span := trace.StartSpan(ctx, name)
	return context.WithValue(newctx, TraceSpanKey{}, span)
}
func (t *tracer) EndSpan(ctx context.Context, httpStatusCode int, err error) {
	ctx.Value(TraceSpanKey{}).(*trace.Span).End()
}

var Tracer = &tracer{}
var RootSpan *trace.Span
