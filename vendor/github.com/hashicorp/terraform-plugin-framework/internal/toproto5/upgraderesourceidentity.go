// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// UpgradeResourceIdentityResponse returns the *tfprotov5.UpgradeResourceIdentityResponse
// equivalent of a *fwserver.UpgradeResourceIdentityResponse.
func UpgradeResourceIdentityResponse(ctx context.Context, fw *fwserver.UpgradeResourceIdentityResponse) *tfprotov5.UpgradeResourceIdentityResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.UpgradeResourceIdentityResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	upgradedIdentity, diags := ResourceIdentity(ctx, fw.UpgradedIdentity)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.UpgradedIdentity = upgradedIdentity

	return proto5
}
