// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto6server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto6"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateProviderConfig satisfies the tfprotov6.ProviderServer interface.
func (s *Server) ValidateProviderConfig(ctx context.Context, proto6Req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ValidateProviderConfigResponse{}

	providerSchema, diags := s.FrameworkServer.ProviderSchema(ctx)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateProviderConfigResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto6.ValidateProviderConfigRequest(ctx, proto6Req, providerSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto6.ValidateProviderConfigResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ValidateProviderConfig(ctx, fwReq, fwResp)

	return toproto6.ValidateProviderConfigResponse(ctx, fwResp), nil
}
