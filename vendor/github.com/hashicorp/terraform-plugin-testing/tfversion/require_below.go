// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// RequireBelow will fail the test if the Terraform CLI
// version is inclusively above the given version. For example, if given
// version.Must(version.NewVersion("1.8.0")), then versions 1.8.x and
// above will fail the test.
//
// Prereleases of Terraform CLI (whether alpha, beta, or rc) are considered
// equal to a given patch version. For example, if given
// version.Must(version.NewVersion("1.8.0")), then 1.8.0-rc1 will fail, not run,
// the test. Terraform prereleases are considered as potential candidates for
// the upcoming version and therefore are treated as semantically equal for
// testing purposes. If failing prereleases of the same patch release is
// desired, give a lower prerelease version. For example, if given
// version.Must(version.NewVersion("1.8.0-rc1")), then 1.8.0-rc2 will fail the
// test.
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
	var terraformVersion *version.Version

	// If given a prerelease version, check the Terraform CLI version directly,
	// otherwise use the core version so that prereleases are treated as equal.
	if s.maximumVersion.Prerelease() != "" {
		terraformVersion = req.TerraformVersion
	} else {
		terraformVersion = req.TerraformVersion.Core()
	}

	if terraformVersion.GreaterThanOrEqual(s.maximumVersion) {
		resp.Error = fmt.Errorf("expected Terraform CLI version below %s but detected version is %s",
			s.maximumVersion, req.TerraformVersion)
	}
}
