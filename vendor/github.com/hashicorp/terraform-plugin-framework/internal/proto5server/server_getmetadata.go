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

// GetMetadata satisfies the tfprotov5.ProviderServer interface.
func (s *Server) GetMetadata(ctx context.Context, proto6Req *tfprotov5.GetMetadataRequest) (*tfprotov5.GetMetadataResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwReq := fromproto5.GetMetadataRequest(ctx, proto6Req)
	fwResp := &fwserver.GetMetadataResponse{}

	s.FrameworkServer.GetMetadata(ctx, fwReq, fwResp)

	return toproto5.GetMetadataResponse(ctx, fwResp), nil
}
