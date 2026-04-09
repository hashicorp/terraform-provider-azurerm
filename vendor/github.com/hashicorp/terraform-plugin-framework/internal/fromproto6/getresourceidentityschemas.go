// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// GetResourceIdentitySchemasRequest returns the *fwserver.GetResourceIdentitySchemasRequest
// equivalent of a *tfprotov6.GetResourceIdentitySchemasRequest.
func GetResourceIdentitySchemasRequest(ctx context.Context, proto6 *tfprotov6.GetResourceIdentitySchemasRequest) *fwserver.GetResourceIdentitySchemasRequest {
	if proto6 == nil {
		return nil
	}

	fw := &fwserver.GetResourceIdentitySchemasRequest{}

	return fw
}
