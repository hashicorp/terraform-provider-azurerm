// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// FunctionError converts the function error into the tfprotov6 function error.
func FunctionError(ctx context.Context, funcErr *function.FuncError) *tfprotov6.FunctionError {
	if funcErr == nil {
		return nil
	}

	return &tfprotov6.FunctionError{
		Text:             funcErr.Text,
		FunctionArgument: funcErr.FunctionArgument,
	}
}
