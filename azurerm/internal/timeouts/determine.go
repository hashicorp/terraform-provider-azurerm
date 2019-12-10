package timeouts

import (
	"context"
	"fmt"
	"time"

	opencensusTrace "go.opencensus.io/trace"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/tracer"
)

// TODO: tests for this

// ForCreate returns the context wrapped with the timeout for an Create operation
//
// If the 'SupportsCustomTimeouts' feature toggle is enabled - this is wrapped with a context
// Otherwise this returns the default context
func ForCreate(ctx context.Context, d *schema.ResourceData) (context.Context, context.CancelFunc) {
	return buildWithTimeout(ctx, d.Timeout(schema.TimeoutCreate), d.Get("name").(string), "create")
}

// ForCreateUpdate returns the context wrapped with the timeout for an combined Create/Update operation
//
// If the 'SupportsCustomTimeouts' feature toggle is enabled - this is wrapped with a context
// Otherwise this returns the default context
func ForCreateUpdate(ctx context.Context, d *schema.ResourceData) (context.Context, context.CancelFunc) {
	if d.IsNewResource() {
		return ForCreate(ctx, d)
	}

	return ForUpdate(ctx, d)
}

// ForDelete returns the context wrapped with the timeout for an Delete operation
//
// If the 'SupportsCustomTimeouts' feature toggle is enabled - this is wrapped with a context
// Otherwise this returns the default context
func ForDelete(ctx context.Context, d *schema.ResourceData) (context.Context, context.CancelFunc) {
	return buildWithTimeout(ctx, d.Timeout(schema.TimeoutDelete), d.Get("name").(string), "delete")
}

// ForRead returns the context wrapped with the timeout for an Read operation
//
// If the 'SupportsCustomTimeouts' feature toggle is enabled - this is wrapped with a context
// Otherwise this returns the default context
func ForRead(ctx context.Context, d *schema.ResourceData) (context.Context, context.CancelFunc) {
	return buildWithTimeout(ctx, d.Timeout(schema.TimeoutRead), d.Get("name").(string), "read")
}

// ForUpdate returns the context wrapped with the timeout for an Update operation
//
// If the 'SupportsCustomTimeouts' feature toggle is enabled - this is wrapped with a context
// Otherwise this returns the default context
func ForUpdate(ctx context.Context, d *schema.ResourceData) (context.Context, context.CancelFunc) {
	return buildWithTimeout(ctx, d.Timeout(schema.TimeoutUpdate), d.Get("name").(string), "update")
}

func buildWithTimeout(ctx context.Context, timeout time.Duration, name, opname string) (context.Context, context.CancelFunc) {
	if features.SupportsCustomTimeouts() {
		return context.WithTimeout(ctx, timeout)
	}
	nullFunc := func() {
		// do nothing on cancel since timeouts aren't enabled
	}

	if tracer.TracingEnabled() {
		var span *opencensusTrace.Span
		ctx, span = opencensusTrace.StartSpanWithRemoteParent(ctx, fmt.Sprintf("%s: %s", name, opname), tracer.RootSpan.SpanContext())
		originNullFunc := nullFunc
		nullFunc = func() {
			originNullFunc()
			span.End()
		}
	}

	return ctx, nullFunc
}
