// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// CallFunctionRequest returns the *fwserver.CallFunctionRequest
// equivalent of a *tfprotov5.CallFunctionRequest.
func CallFunctionRequest(ctx context.Context, proto *tfprotov5.CallFunctionRequest, function function.Function, functionDefinition function.Definition) (*fwserver.CallFunctionRequest, *function.FuncError) {
	if proto == nil {
		return nil, nil
	}

	fw := &fwserver.CallFunctionRequest{
		Function:           function,
		FunctionDefinition: functionDefinition,
	}

	arguments, funcError := ArgumentsData(ctx, proto.Arguments, functionDefinition)

	fw.Arguments = arguments

	return fw, funcError
}
