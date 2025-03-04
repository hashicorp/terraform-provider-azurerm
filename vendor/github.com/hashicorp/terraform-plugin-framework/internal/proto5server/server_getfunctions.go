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

// GetFunctions satisfies the tfprotov5.ProviderServer interface.
func (s *Server) GetFunctions(ctx context.Context, protoReq *tfprotov5.GetFunctionsRequest) (*tfprotov5.GetFunctionsResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwReq := fromproto5.GetFunctionsRequest(ctx, protoReq)
	fwResp := &fwserver.GetFunctionsResponse{}

	s.FrameworkServer.GetFunctions(ctx, fwReq, fwResp)

	return toproto5.GetFunctionsResponse(ctx, fwResp), nil
}
