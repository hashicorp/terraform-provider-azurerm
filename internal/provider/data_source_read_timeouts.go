// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

// wrapDataSourceReadTimeouts works around https://github.com/hashicorp/terraform-plugin-sdk/issues/1038:
// the Plugin SDK never attaches the schema-declared Timeouts (or the user's `timeouts` block) to a Data
// Source read, so `d.Timeout(TimeoutRead)` always returns the SDK-global default of 20 minutes rather
// than the value declared on the Data Source or configured by the user.
//
// This wraps the Data Source's read function to determine the effective timeout (the user-configured
// `timeouts` block when set, else the schema-declared read timeout) and:
//
//   - registers it with the timeouts package for the duration of the read, so that the
//     `timeouts.ForRead` call inside legacy (context-less) read functions picks it up
//   - additionally applies it as a context deadline for context-aware read functions
//     (i.e. typed Data Sources built via sdk.DataSourceWrapper)
func wrapDataSourceReadTimeouts(r *schema.Resource) {
	if r.Timeouts == nil || r.Timeouts.Read == nil {
		return
	}
	declared := *r.Timeouts.Read

	if r.ReadContext != nil {
		inner := r.ReadContext
		r.ReadContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			timeout := effectiveReadTimeout(d.GetRawConfig(), declared)
			timeouts.RegisterDataSourceReadTimeout(d, timeout)
			defer timeouts.DeregisterDataSourceReadTimeout(d)

			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			return inner(ctx, d, meta)
		}
		return
	}

	if r.Read != nil {
		inner := r.Read
		r.Read = func(d *schema.ResourceData, meta interface{}) error {
			timeouts.RegisterDataSourceReadTimeout(d, effectiveReadTimeout(d.GetRawConfig(), declared))
			defer timeouts.DeregisterDataSourceReadTimeout(d)
			return inner(d, meta)
		}
	}
}

// effectiveReadTimeout returns the read timeout from the `timeouts` block within the raw config
// when one is set, falling back to the schema-declared read timeout.
func effectiveReadTimeout(rawConfig cty.Value, declared time.Duration) time.Duration {
	if rawConfig.IsNull() || !rawConfig.Type().IsObjectType() || !rawConfig.Type().HasAttribute("timeouts") {
		return declared
	}
	block := rawConfig.GetAttr("timeouts")
	if block.IsNull() || !block.IsKnown() || !block.Type().IsObjectType() || !block.Type().HasAttribute("read") {
		return declared
	}
	value := block.GetAttr("read")
	if value.IsNull() || !value.IsKnown() || value.Type() != cty.String {
		return declared
	}
	timeout, err := time.ParseDuration(value.AsString())
	if err != nil {
		return declared
	}
	return timeout
}
