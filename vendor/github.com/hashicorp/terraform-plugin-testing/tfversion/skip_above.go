// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// SkipAbove will skip (pass) the test if the Terraform CLI
// version is exclusively above the given version. For example, if given
// version.Must(version.NewVersion("1.8.0")), then 1.8.x or
// any other later versions will skip the test.
//
// Prereleases of Terraform CLI (whether alpha, beta, or rc) are considered
// equal to a given patch version. For example, if given
// version.Must(version.NewVersion("1.8.0")), then 1.8.0-rc1 will run, not skip,
// the test. Terraform prereleases are considered as potential candidates for
// the upcoming version and therefore are treated as semantically equal for
// testing. If skipping prereleases of the same patch release is desired, give a
// lower prerelease version. For example, if given
// version.Must(version.NewVersion("1.8.0-rc1")), then 1.8.0-rc2 will skip the
// test.
func SkipAbove(maximumVersion *version.Version) TerraformVersionCheck {
	return skipAboveCheck{
		maximumVersion: maximumVersion,
	}
}

// skipAboveCheck implements the TerraformVersionCheck interface
type skipAboveCheck struct {
	maximumVersion *version.Version
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (s skipAboveCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {
	var terraformVersion *version.Version

	// If given a prerelease version, check the Terraform CLI version directly,
	// otherwise use the core version so that prereleases are treated as equal.
	if s.maximumVersion.Prerelease() != "" {
		terraformVersion = req.TerraformVersion
	} else {
		terraformVersion = req.TerraformVersion.Core()
	}

	if terraformVersion.GreaterThan(s.maximumVersion) {
		resp.Skip = fmt.Sprintf("Terraform CLI version %s is above maximum version %s: skipping test",
			req.TerraformVersion, s.maximumVersion)
	}
}
