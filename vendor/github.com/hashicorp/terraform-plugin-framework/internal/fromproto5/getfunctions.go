// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// GetFunctionsRequest returns the *fwserver.GetFunctionsRequest
// equivalent of a *tfprotov5.GetFunctionsRequest.
func GetFunctionsRequest(ctx context.Context, proto *tfprotov5.GetFunctionsRequest) *fwserver.GetFunctionsRequest {
	if proto == nil {
		return nil
	}

	fw := &fwserver.GetFunctionsRequest{}

	return fw
}
