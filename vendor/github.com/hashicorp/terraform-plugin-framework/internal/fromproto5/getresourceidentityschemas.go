// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// GetResourceIdentitySchemasRequest returns the *fwserver.GetResourceIdentitySchemasRequest
// equivalent of a *tfprotov5.GetResourceIdentitySchemasRequest.
func GetResourceIdentitySchemasRequest(ctx context.Context, proto5 *tfprotov5.GetResourceIdentitySchemasRequest) *fwserver.GetResourceIdentitySchemasRequest {
	if proto5 == nil {
		return nil
	}

	fw := &fwserver.GetResourceIdentitySchemasRequest{}

	return fw
}
