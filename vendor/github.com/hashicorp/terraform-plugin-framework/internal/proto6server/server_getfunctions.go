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

// GetFunctions satisfies the tfprotov6.ProviderServer interface.
func (s *Server) GetFunctions(ctx context.Context, protoReq *tfprotov6.GetFunctionsRequest) (*tfprotov6.GetFunctionsResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwReq := fromproto6.GetFunctionsRequest(ctx, protoReq)
	fwResp := &fwserver.GetFunctionsResponse{}

	s.FrameworkServer.GetFunctions(ctx, fwReq, fwResp)

	return toproto6.GetFunctionsResponse(ctx, fwResp), nil
}
