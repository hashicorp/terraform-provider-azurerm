// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// GetFunctionsResponse returns the *tfprotov6.GetFunctionsResponse
// equivalent of a *fwserver.GetFunctionsResponse.
func GetFunctionsResponse(ctx context.Context, fw *fwserver.GetFunctionsResponse) *tfprotov6.GetFunctionsResponse {
	if fw == nil {
		return nil
	}

	proto := &tfprotov6.GetFunctionsResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
		Functions:   make(map[string]*tfprotov6.Function, len(fw.FunctionDefinitions)),
	}

	for name, functionDefinition := range fw.FunctionDefinitions {
		proto.Functions[name] = Function(ctx, functionDefinition)
	}

	return proto
}
