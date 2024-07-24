// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto5server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto5"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto5"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ConfigureProvider satisfies the tfprotov5.ProviderServer interface.
func (s *Server) ConfigureProvider(ctx context.Context, proto5Req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &provider.ConfigureResponse{}

	providerSchema, diags := s.FrameworkServer.ProviderSchema(ctx)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ConfigureProviderResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.ConfigureProviderRequest(ctx, proto5Req, providerSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.ConfigureProviderResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ConfigureProvider(ctx, fwReq, fwResp)

	return toproto5.ConfigureProviderResponse(ctx, fwResp), nil
}
