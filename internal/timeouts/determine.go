// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package timeouts

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// ForCreate returns the context wrapped with the timeout for an Create operation
//
// If the 'SupportsCustomTimeouts' feature toggle is enabled - this is wrapped with a context
// Otherwise this returns the default context
func ForCreate(ctx context.Context, d *pluginsdk.ResourceData) (context.Context, context.CancelFunc) {
	return buildWithTimeout(ctx, d.Timeout(pluginsdk.TimeoutCreate))
}

// ForCreateUpdate returns the context wrapped with the timeout for an combined Create/Update operation
//
// If the 'SupportsCustomTimeouts' feature toggle is enabled - this is wrapped with a context
// Otherwise this returns the default context
func ForCreateUpdate(ctx context.Context, d *pluginsdk.ResourceData) (context.Context, context.CancelFunc) {
	if d.IsNewResource() {
		return ForCreate(ctx, d)
	}

	return ForUpdate(ctx, d)
}

// ForDelete returns the context wrapped with the timeout for an Delete operation
//
// If the 'SupportsCustomTimeouts' feature toggle is enabled - this is wrapped with a context
// Otherwise this returns the default context
func ForDelete(ctx context.Context, d *pluginsdk.ResourceData) (context.Context, context.CancelFunc) {
	return buildWithTimeout(ctx, d.Timeout(pluginsdk.TimeoutDelete))
}

// ForRead returns the context wrapped with the timeout for an Read operation
//
// If the 'SupportsCustomTimeouts' feature toggle is enabled - this is wrapped with a context
// Otherwise this returns the default context
func ForRead(ctx context.Context, d *pluginsdk.ResourceData) (context.Context, context.CancelFunc) {
	// for Data Source reads `d.Timeout` unconditionally returns the SDK-global default of 20
	// minutes (hashicorp/terraform-plugin-sdk#1038) - so when the provider has registered the
	// effective timeout for this read (see internal/provider) prefer that value
	if timeout, ok := dataSourceReadTimeout(d); ok {
		return buildWithTimeout(ctx, timeout)
	}
	return buildWithTimeout(ctx, d.Timeout(pluginsdk.TimeoutRead))
}

// ForUpdate returns the context wrapped with the timeout for an Update operation
//
// If the 'SupportsCustomTimeouts' feature toggle is enabled - this is wrapped with a context
// Otherwise this returns the default context
func ForUpdate(ctx context.Context, d *pluginsdk.ResourceData) (context.Context, context.CancelFunc) {
	return buildWithTimeout(ctx, d.Timeout(pluginsdk.TimeoutUpdate))
}

func buildWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}
