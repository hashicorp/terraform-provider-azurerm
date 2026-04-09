// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// GetProviderSchemaRequest returns the *fwserver.GetProviderSchemaRequest
// equivalent of a *tfprotov5.GetProviderSchemaRequest.
func GetProviderSchemaRequest(ctx context.Context, proto5 *tfprotov5.GetProviderSchemaRequest) *fwserver.GetProviderSchemaRequest {
	if proto5 == nil {
		return nil
	}

	fw := &fwserver.GetProviderSchemaRequest{}

	return fw
}
