// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// GetMetadataRequest returns the *fwserver.GetMetadataRequest
// equivalent of a *tfprotov6.GetMetadataRequest.
func GetMetadataRequest(ctx context.Context, proto6 *tfprotov6.GetMetadataRequest) *fwserver.GetMetadataRequest {
	if proto6 == nil {
		return nil
	}

	fw := &fwserver.GetMetadataRequest{}

	return fw
}
