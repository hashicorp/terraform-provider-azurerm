// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// RequireAbove will fail the test if the Terraform CLI
// version is exclusively below the given version. For example, if given
// version.Must(version.NewVersion("1.8.0")), then 1.7.x or
// any other prior versions will fail the test.
//
// Prereleases of Terraform CLI (whether alpha, beta, or rc) are considered
// equal to a given patch version. For example, if given
// version.Must(version.NewVersion("1.8.0")), then 1.8.0-rc1 will run, not fail,
// the test. Terraform prereleases are considered as potential candidates for
// the upcoming version and therefore are treated as semantically equal for
// testing. If failing prereleases of the same patch release is desired, give a
// higher prerelease version. For example, if given
// version.Must(version.NewVersion("1.8.0-rc2")), then 1.8.0-rc1 will fail the
// test.
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
	var terraformVersion *version.Version

	// If given a prerelease version, check the Terraform CLI version directly,
	// otherwise use the core version so that prereleases are treated as equal.
	if r.minimumVersion.Prerelease() != "" {
		terraformVersion = req.TerraformVersion
	} else {
		terraformVersion = req.TerraformVersion.Core()
	}

	if terraformVersion.LessThan(r.minimumVersion) {
		resp.Error = fmt.Errorf("expected Terraform CLI version above %s but detected version is %s",
			r.minimumVersion, req.TerraformVersion)
	}
}
