// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// RequireNot will fail the test if the Terraform CLI
// version matches the given version.
//
// Prereleases of Terraform CLI (whether alpha, beta, or rc) are considered
// equal to a given patch version. For example, if given
// version.Must(version.NewVersion("1.8.0")), then 1.8.0-rc1 will fail, not run,
// the test. Terraform prereleases are considered as potential candidates for
// the upcoming version and therefore are treated as semantically equal for
// testing purposes. If running prereleases of the same patch release is
// desired, give a different prerelease version. For example, if given
// version.Must(version.NewVersion("1.8.0-rc2")), then 1.8.0-rc1 will
// run the test.
func RequireNot(version *version.Version) TerraformVersionCheck {
	return requireNotCheck{
		version: version,
	}
}

// requireNotCheck implements the TerraformVersionCheck interface
type requireNotCheck struct {
	version *version.Version
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (s requireNotCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {
	var terraformVersion *version.Version

	// If given a prerelease version, check the Terraform CLI version directly,
	// otherwise use the core version so that prereleases are treated as equal.
	if s.version.Prerelease() != "" {
		terraformVersion = req.TerraformVersion
	} else {
		terraformVersion = req.TerraformVersion.Core()
	}

	if terraformVersion.Equal(s.version) {
		resp.Error = fmt.Errorf("unexpected Terraform CLI version: %s", s.version)
	}
}
