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

// GetResourceIdentitySchemas satisfies the tfprotov5.ProviderServer interface.
func (s *Server) GetResourceIdentitySchemas(ctx context.Context, proto5Req *tfprotov5.GetResourceIdentitySchemasRequest) (*tfprotov5.GetResourceIdentitySchemasResponse, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

	fwReq := fromproto5.GetResourceIdentitySchemasRequest(ctx, proto5Req)
	fwResp := &fwserver.GetResourceIdentitySchemasResponse{}

	s.FrameworkServer.GetResourceIdentitySchemas(ctx, fwReq, fwResp)

	return toproto5.GetResourceIdentitySchemasResponse(ctx, fwResp), nil
}
