// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// GetMetadataRequest returns the *fwserver.GetMetadataRequest
// equivalent of a *tfprotov5.GetMetadataRequest.
func GetMetadataRequest(ctx context.Context, proto5 *tfprotov5.GetMetadataRequest) *fwserver.GetMetadataRequest {
	if proto5 == nil {
		return nil
	}

	fw := &fwserver.GetMetadataRequest{}

	return fw
}
