// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto5server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto5"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// PrepareProviderConfig satisfies the tfprotov5.ProviderServer interface.
func (s *Server) PrepareProviderConfig(ctx context.Context, proto5Req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwResp := &fwserver.ValidateProviderConfigResponse{}

	providerSchema, diags := s.FrameworkServer.ProviderSchema(ctx)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.PrepareProviderConfigResponse(ctx, fwResp), nil
	}

	fwReq, diags := fromproto5.PrepareProviderConfigRequest(ctx, proto5Req, providerSchema)

	fwResp.Diagnostics.Append(diags...)

	if fwResp.Diagnostics.HasError() {
		return toproto5.PrepareProviderConfigResponse(ctx, fwResp), nil
	}

	s.FrameworkServer.ValidateProviderConfig(ctx, fwReq, fwResp)

	return toproto5.PrepareProviderConfigResponse(ctx, fwResp), nil
}
