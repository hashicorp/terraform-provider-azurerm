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

// GetProviderSchema satisfies the tfprotov6.ProviderServer interface.
func (s *Server) GetProviderSchema(ctx context.Context, proto6Req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwReq := fromproto6.GetProviderSchemaRequest(ctx, proto6Req)
	fwResp := &fwserver.GetProviderSchemaResponse{}

	s.FrameworkServer.GetProviderSchema(ctx, fwReq, fwResp)

	return toproto6.GetProviderSchemaResponse(ctx, fwResp), nil
}
