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
	// 1. Configure exporter to export traces to Zipkin.
	localEndpoint, err := openzipkin.NewEndpoint("terraform", "192.168.1.5:5454")
	if err != nil {
		log.Fatalf("Failed to create the local zipkinEndpoint: %v", err)
	}
	reporter := zipkinHTTP.NewReporter("http://localhost:9411/api/v2/spans")
	ze := zipkin.NewExporter(reporter, localEndpoint)
	trace.RegisterExporter(ze)

	// 2. Configure 100% sample rate, otherwise, few traces will be sampled.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
}

const MyKey = "foo"

type myTracer struct{}

func (t *myTracer) NewTransport(base *http.Transport) http.RoundTripper {
	return base
}

func (t *myTracer) StartSpan(ctx context.Context, name string) context.Context {
	newctx, span := trace.StartSpan(ctx, name)
	return context.WithValue(newctx, MyKey, span)
}
func (t *myTracer) EndSpan(ctx context.Context, httpStatusCode int, err error) {
	ctx.Value(MyKey).(*trace.Span).End()
}

var MyTracer = &myTracer{}
