// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto6server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto6"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto6"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ConfigureProvider satisfies the tfprotov6.ProviderServer interface.
func (s *Server) ConfigureProvider(ctx context.Context, proto6Req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &provider.ConfigureResponse{}

	providerSchema, diags := s.FrameworkServer.ProviderSchema(ctx)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ConfigureProviderResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.ConfigureProviderRequest(ctx, proto6Req, providerSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ConfigureProviderResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ConfigureProvider(ctx, fwReq, fwResp)

	return toproto6.ConfigureProviderResponse(ctx, fwResp), nil
}
