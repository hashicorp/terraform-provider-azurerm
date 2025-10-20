// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// SkipBelow will skip (pass) the test if the Terraform CLI
// version is exclusively below the given version. For example, if given
// version.Must(version.NewVersion("1.8.0")), then 1.7.x or
// any other prior versions will skip the test.
//
// Prereleases of Terraform CLI (whether alpha, beta, or rc) are considered
// equal to a given patch version. For example, if given
// version.Must(version.NewVersion("1.8.0")), then 1.8.0-rc1 will run, not skip,
// the test. Terraform prereleases are considered as potential candidates for
// the upcoming version and therefore are treated as important for testing to
// run. If skipping prereleases of the same patch release is desired, give a
// higher prerelease version. For example, if given
// version.Must(version.NewVersion("1.8.0-rc2")), then 1.8.0-rc1 will skip the
// test.
func SkipBelow(minimumVersion *version.Version) TerraformVersionCheck {
	return skipBelowCheck{
		minimumVersion: minimumVersion,
	}
}

// skipBelowCheck implements the TerraformVersionCheck interface
type skipBelowCheck struct {
	minimumVersion *version.Version
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (s skipBelowCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {
	var terraformVersion *version.Version

	// If given a prerelease version, check the Terraform CLI version directly,
	// otherwise use the core version so that prereleases are treated as equal.
	if s.minimumVersion.Prerelease() != "" {
		terraformVersion = req.TerraformVersion
	} else {
		terraformVersion = req.TerraformVersion.Core()
	}

	if terraformVersion.LessThan(s.minimumVersion) {
		resp.Skip = fmt.Sprintf("Terraform CLI version %s is below minimum version %s: skipping test",
			req.TerraformVersion, s.minimumVersion)
	}
}
