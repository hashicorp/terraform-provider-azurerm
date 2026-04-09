// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateDataSourceConfigResponse returns the *tfprotov6.ValidateDataSourceConfigResponse
// equivalent of a *fwserver.ValidateDataSourceConfigResponse.
func ValidateDataSourceConfigResponse(ctx context.Context, fw *fwserver.ValidateDataSourceConfigResponse) *tfprotov6.ValidateDataResourceConfigResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.ValidateDataResourceConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto6
}
