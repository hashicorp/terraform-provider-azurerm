// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// GetProviderSchemaRequest returns the *fwserver.GetProviderSchemaRequest
// equivalent of a *tfprotov6.GetProviderSchemaRequest.
func GetProviderSchemaRequest(ctx context.Context, proto6 *tfprotov6.GetProviderSchemaRequest) *fwserver.GetProviderSchemaRequest {
	if proto6 == nil {
		return nil
	}

	fw := &fwserver.GetProviderSchemaRequest{}

	return fw
}
