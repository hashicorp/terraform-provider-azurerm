// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ValidateDataSourceConfigResponse returns the *tfprotov5.ValidateDataSourceConfigResponse
// equivalent of a *fwserver.ValidateDataSourceConfigResponse.
func ValidateDataSourceConfigResponse(ctx context.Context, fw *fwserver.ValidateDataSourceConfigResponse) *tfprotov5.ValidateDataSourceConfigResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.ValidateDataSourceConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto5
}
