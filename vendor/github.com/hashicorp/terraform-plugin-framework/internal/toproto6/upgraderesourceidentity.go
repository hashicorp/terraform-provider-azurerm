// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// UpgradeResourceIdentityResponse returns the *tfprotov6.UpgradeResourceIdentityResponse
// equivalent of a *fwserver.UpgradeResourceIdentityResponse.
func UpgradeResourceIdentityResponse(ctx context.Context, fw *fwserver.UpgradeResourceIdentityResponse) *tfprotov6.UpgradeResourceIdentityResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.UpgradeResourceIdentityResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	upgradedIdentity, diags := ResourceIdentity(ctx, fw.UpgradedIdentity)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.UpgradedIdentity = upgradedIdentity

	return proto6
}
