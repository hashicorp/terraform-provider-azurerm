// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// ConfigureProvider implements the framework server ConfigureProvider RPC.
func (s *Server) ConfigureProvider(ctx context.Context, req *provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	logging.FrameworkTrace(ctx, "Calling provider defined Provider Configure")

	if req != nil {
		s.Provider.Configure(ctx, *req, resp)
	} else {
		s.Provider.Configure(ctx, provider.ConfigureRequest{}, resp)
	}

	logging.FrameworkTrace(ctx, "Called provider defined Provider Configure")

	if resp.Deferred != nil {
		if !req.ClientCapabilities.DeferralAllowed {
			resp.Diagnostics.AddError("Invalid Deferred Provider Response",
				"Provider configured a deferred response for all resources and data sources but the Terraform request "+
					"did not indicate support for deferred actions. This is an issue with the provider and should be reported to the provider developers.")
			return
		}

		logging.FrameworkDebug(ctx, "Provider has configured a deferred response, "+
			"all associated resources and data sources will automatically return a deferred response.")
	}

	s.deferred = resp.Deferred
	s.DataSourceConfigureData = resp.DataSourceData
	s.ResourceConfigureData = resp.ResourceData
	s.EphemeralResourceConfigureData = resp.EphemeralResourceData
	s.ActionConfigureData = resp.ActionData
	s.ListResourceConfigureData = resp.ListResourceData
}
