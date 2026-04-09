// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// GetFunctionsResponse returns the *tfprotov5.GetFunctionsResponse
// equivalent of a *fwserver.GetFunctionsResponse.
func GetFunctionsResponse(ctx context.Context, fw *fwserver.GetFunctionsResponse) *tfprotov5.GetFunctionsResponse {
	if fw == nil {
		return nil
	}

	proto := &tfprotov5.GetFunctionsResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
		Functions:   make(map[string]*tfprotov5.Function, len(fw.FunctionDefinitions)),
	}

	for name, functionDefinition := range fw.FunctionDefinitions {
		proto.Functions[name] = Function(ctx, functionDefinition)
	}

	return proto
}
