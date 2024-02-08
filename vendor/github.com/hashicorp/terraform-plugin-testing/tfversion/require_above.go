// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// RequireAbove will fail the test if the Terraform CLI
// version is below the given version. For example, if given
// version.Must(version.NewVersion("0.15.0")), then 0.14.x or
// any other prior minor versions will fail the test.
func RequireAbove(minimumVersion *version.Version) TerraformVersionCheck {
	return requireAboveCheck{
		minimumVersion: minimumVersion,
	}
}

// requireAboveCheck implements the TerraformVersionCheck interface
type requireAboveCheck struct {
	minimumVersion *version.Version
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (r requireAboveCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {

	if req.TerraformVersion.LessThan(r.minimumVersion) {
		resp.Error = fmt.Errorf("expected Terraform CLI version above %s but detected version is %s",
			r.minimumVersion, req.TerraformVersion)
	}
}
