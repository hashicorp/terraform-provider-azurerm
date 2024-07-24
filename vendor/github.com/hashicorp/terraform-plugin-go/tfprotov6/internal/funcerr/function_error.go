// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcerr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/internal/logging"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// FunctionError is a single FunctionError.
type FunctionError tfprotov6.FunctionError

// HasError returns true if the FunctionError is not empty.
func (e *FunctionError) HasError() bool {
	if e == nil {
		return false
	}

	return e.Text != "" || e.FunctionArgument != nil
}

// Log will log the function error:
func (e *FunctionError) Log(ctx context.Context) {
	if e == nil {
		return
	}

	if !e.HasError() {
		return
	}

	switch {
	case e.FunctionArgument != nil && e.Text != "":
		logging.ProtocolError(ctx, "Response contains function error", map[string]interface{}{
			logging.KeyFunctionErrorText:     e.Text,
			logging.KeyFunctionErrorArgument: *e.FunctionArgument,
		})
	case e.FunctionArgument != nil:
		logging.ProtocolError(ctx, "Response contains function error", map[string]interface{}{
			logging.KeyFunctionErrorArgument: *e.FunctionArgument,
		})
	case e.Text != "":
		logging.ProtocolError(ctx, "Response contains function error", map[string]interface{}{
			logging.KeyFunctionErrorText: e.Text,
		})
	}
}
