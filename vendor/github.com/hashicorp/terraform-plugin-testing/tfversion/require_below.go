// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// RequireBelow will fail the test if the Terraform CLI
// version is above the given version. For example, if given
// version.Must(version.NewVersion("0.15.0")), then versions 0.15.x and
// above will fail the test.
func RequireBelow(maximumVersion *version.Version) TerraformVersionCheck {
	return requireBelowCheck{
		maximumVersion: maximumVersion,
	}
}

// requireBelowCheck implements the TerraformVersionCheck interface
type requireBelowCheck struct {
	maximumVersion *version.Version
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (s requireBelowCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {

	if req.TerraformVersion.GreaterThan(s.maximumVersion) {
		resp.Error = fmt.Errorf("expected Terraform CLI version below %s but detected version is %s",
			s.maximumVersion, req.TerraformVersion)
	}
}
