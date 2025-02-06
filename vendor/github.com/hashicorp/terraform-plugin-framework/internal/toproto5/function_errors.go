// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// FunctionError converts the function error into the tfprotov5 function error.
func FunctionError(ctx context.Context, funcErr *function.FuncError) *tfprotov5.FunctionError {
	if funcErr == nil {
		return nil
	}

	return &tfprotov5.FunctionError{
		Text:             funcErr.Text,
		FunctionArgument: funcErr.FunctionArgument,
	}
}
