// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// SkipBetween will skip the test if the Terraform CLI
// version is between the given minimum (inclusive) and maximum (exclusive).
// For example, if given a minimum version of version.Must(version.NewVersion("1.7.0"))
// and a maximum version of version.Must(version.NewVersion("1.8.0")), then versions 1.7.x
// will skip the test.
//
// Prereleases of Terraform CLI (whether alpha, beta, or rc) are considered
// equal to a given patch version. For example, if given a minimum version of
// version.Must(version.NewVersion("1.8.0")), then 1.8.0-rc1 will skip, not run,
// the test. Terraform prereleases are considered as potential candidates for
// the upcoming version and therefore are treated as semantically equal for
// testing purposes. If running prereleases of the same patch release is
// desired, give a higher prerelease version. For example, if given a minimum
// version of version.Must(version.NewVersion("1.8.0-rc2")), then 1.8.0-rc1 will
// run the test.
func SkipBetween(minimumVersion, maximumVersion *version.Version) TerraformVersionCheck {
	return skipBetweenCheck{
		minimumVersion: minimumVersion,
		maximumVersion: maximumVersion,
	}
}

// skipBetweenCheck implements the TerraformVersionCheck interface
type skipBetweenCheck struct {
	minimumVersion *version.Version
	maximumVersion *version.Version
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (s skipBetweenCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {
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

	if minTerraformVersion.GreaterThanOrEqual(s.minimumVersion) && maxTerraformVersion.LessThan(s.maximumVersion) {
		resp.Skip = fmt.Sprintf("Terraform CLI version %s is between %s and %s: skipping test.",
			req.TerraformVersion, s.minimumVersion, s.maximumVersion)
	}
}
