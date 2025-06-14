// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
)

// GetResourceIdentitySchemasRequest is the framework server request for the
// GetResourceIdentitySchemas RPC.
type GetResourceIdentitySchemasRequest struct{}

// GetResourceIdentitySchemasResponse is the framework server response for the
// GetResourceIdentitySchemas RPC.
type GetResourceIdentitySchemasResponse struct {
	IdentitySchemas map[string]fwschema.Schema
	Diagnostics     diag.Diagnostics
}

// GetResourceIdentitySchemas implements the framework server GetResourceIdentitySchemas RPC.
func (s *Server) GetResourceIdentitySchemas(ctx context.Context, req *GetResourceIdentitySchemasRequest, resp *GetResourceIdentitySchemasResponse) {
	identitySchemas, diags := s.ResourceIdentitySchemas(ctx)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.IdentitySchemas = identitySchemas
}
