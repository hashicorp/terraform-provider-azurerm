// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// CallFunctionResponse returns the *tfprotov5.CallFunctionResponse
// equivalent of a *fwserver.CallFunctionResponse.
func CallFunctionResponse(ctx context.Context, fw *fwserver.CallFunctionResponse) *tfprotov5.CallFunctionResponse {
	if fw == nil {
		return nil
	}

	result, resultErr := FunctionResultData(ctx, fw.Result)

	funcErr := function.ConcatFuncErrors(fw.Error, resultErr)

	return &tfprotov5.CallFunctionResponse{
		Error:  FunctionError(ctx, funcErr),
		Result: result,
	}
}
