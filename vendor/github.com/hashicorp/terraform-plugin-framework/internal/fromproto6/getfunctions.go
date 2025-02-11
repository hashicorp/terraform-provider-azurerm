// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// GetFunctionsRequest returns the *fwserver.GetFunctionsRequest
// equivalent of a *tfprotov6.GetFunctionsRequest.
func GetFunctionsRequest(ctx context.Context, proto *tfprotov6.GetFunctionsRequest) *fwserver.GetFunctionsRequest {
	if proto == nil {
		return nil
	}

	fw := &fwserver.GetFunctionsRequest{}

	return fw
}
