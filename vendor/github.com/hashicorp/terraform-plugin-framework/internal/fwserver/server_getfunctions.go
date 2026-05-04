// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

// GetFunctionsRequest is the framework server request for the
// GetFunctions RPC.
type GetFunctionsRequest struct{}

// GetFunctionsResponse is the framework server response for the
// GetFunctions RPC.
type GetFunctionsResponse struct {
	FunctionDefinitions map[string]function.Definition
	Diagnostics         diag.Diagnostics
}

// GetFunctions implements the framework server GetFunctions RPC.
func (s *Server) GetFunctions(ctx context.Context, req *GetFunctionsRequest, resp *GetFunctionsResponse) {
	resp.FunctionDefinitions = map[string]function.Definition{}

	functionDefinitions, diags := s.FunctionDefinitions(ctx)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.FunctionDefinitions = functionDefinitions
}
