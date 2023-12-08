// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// RequireBetween will fail the test if the Terraform CLI
// version is outside the given minimum (exclusive) and maximum (inclusive).
// For example, if given a minimum version of version.Must(version.NewVersion("0.15.0"))
// and a maximum version of version.Must(version.NewVersion("1.0.0")), then 0.15.x or
// any other prior versions and versions greater than 1.0.0 will fail the test.
func RequireBetween(minimumVersion, maximumVersion *version.Version) TerraformVersionCheck {
	return requireBetweenCheck{
		minimumVersion: minimumVersion,
		maximumVersion: maximumVersion,
	}
}

// requireBetweenCheck implements the TerraformVersionCheck interface
type requireBetweenCheck struct {
	minimumVersion *version.Version
	maximumVersion *version.Version
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (s requireBetweenCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {

	if req.TerraformVersion.LessThan(s.minimumVersion) || req.TerraformVersion.GreaterThanOrEqual(s.maximumVersion) {
		resp.Error = fmt.Errorf("expected Terraform CLI version between %s and %s but detected version is %s",
			s.minimumVersion, s.maximumVersion, req.TerraformVersion)
	}
}
