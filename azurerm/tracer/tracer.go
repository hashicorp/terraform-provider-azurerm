package tracer

import (
	"context"
	"log"
	"net/http"
	"os"

	autorestTracing "github.com/Azure/go-autorest/tracing"

	_ "github.com/Azure/go-autorest/tracing/opencensus"

	opencensusTrace "go.opencensus.io/trace"
)

const (
	envTfAzureTracer = "TF_AZURE_TRACER"
)

// RootSpan is the first span created during the lifetime of plugin server.
// Based on this span, children spans will be created on each resource operation.
var RootSpan *opencensusTrace.Span

// TracingEnabled tells whether tracing is enabled
func TracingEnabled() bool {
	return os.Getenv(envTfAzureTracer) != ""
}

// tracer implements Tracer interface in go-autorest
type tracer struct{}

type traceSpanKey struct{}

func (t *tracer) NewTransport(base *http.Transport) http.RoundTripper {
	return base
}

func (t *tracer) StartSpan(ctx context.Context, name string) context.Context {
	newctx, span := opencensusTrace.StartSpan(ctx, name)
	return context.WithValue(newctx, traceSpanKey{}, span)
}
func (t *tracer) EndSpan(ctx context.Context, httpStatusCode int, err error) {
	ctx.Value(traceSpanKey{}).(*opencensusTrace.Span).End()
}

type exporterBuilder func(string) (opencensusTrace.Exporter, error)

var exporterBuilders = map[string]exporterBuilder{
	"jaeger": buildJaegerExporter,
	"zipkin": buildZipkinExporter,
}

func Init() {
	// enable tracing in azure sdk
	if err := os.Setenv("AZURE_SDK_TRACING_ENABLED", "true"); err != nil {
		log.Fatalf("[ERROR] failed set env var: %v", err)
	}

	// build exporter of user's choice
	t := os.Getenv(envTfAzureTracer)
	b, ok := exporterBuilders[t]
	if !ok {
		log.Fatalf("[ERROR] unknown tracer: %s", t)
	}
	exporter, err := b("terraform")
	if err != nil {
		log.Fatalf("[ERROR] failed to build exporter: %v", err)
	}

	// register exporter to opencensus
	opencensusTrace.RegisterExporter(exporter)
	opencensusTrace.ApplyConfig(opencensusTrace.Config{DefaultSampler: opencensusTrace.AlwaysSample()})

	// register tracer to go-autorest
	autorestTracing.Register(&tracer{})
}
