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
// For example, if given a minimum version of version.Must(version.NewVersion("1.7.0"))
// and a maximum version of version.Must(version.NewVersion("1.8.0")), then 1.6.x or
// any other prior versions and versions greater than or equal to 1.8.0 will fail the test.
//
// Prereleases of Terraform CLI (whether alpha, beta, or rc) are considered
// equal to a given patch version. For example, if given a minimum version of
// version.Must(version.NewVersion("1.8.0")), then 1.8.0-rc1 will run, not fail,
// the test. Terraform prereleases are considered as potential candidates for
// the upcoming version and therefore are treated as semantically equal for
// testing purposes. If failing prereleases of the same patch release is
// desired, give a higher prerelease version. For example, if given a minimum
// version of version.Must(version.NewVersion("1.8.0-rc2")), then 1.8.0-rc1 will
// fail the test.
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
	var maxTerraformVersion, minTerraformVersion *version.Version

	// If given a prerelease maximum version, check the Terraform CLI version
	// directly, otherwise use the core version so that prereleases are treated
	// as equal.
	if s.maximumVersion.Prerelease() != "" {
		maxTerraformVersion = req.TerraformVersion
	} else {
		maxTerraformVersion = req.TerraformVersion.Core()
	}

	// If given a prerelease minimum version, check the Terraform CLI version
	// directly, otherwise use the core version so that prereleases are treated
	// as equal.
	if s.minimumVersion.Prerelease() != "" {
		minTerraformVersion = req.TerraformVersion
	} else {
		minTerraformVersion = req.TerraformVersion.Core()
	}

	if minTerraformVersion.LessThan(s.minimumVersion) || maxTerraformVersion.GreaterThanOrEqual(s.maximumVersion) {
		resp.Error = fmt.Errorf("expected Terraform CLI version between %s and %s but detected version is %s",
			s.minimumVersion, s.maximumVersion, req.TerraformVersion)
	}
}
