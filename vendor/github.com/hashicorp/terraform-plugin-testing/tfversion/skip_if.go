// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// SkipIf will skip (pass) the test if the Terraform CLI
// version matches the given version.
//
// Prereleases of Terraform CLI (whether alpha, beta, or rc) are considered
// equal to a given patch version. For example, if given
// version.Must(version.NewVersion("1.8.0")), then 1.8.0-rc1 will skip, not run,
// the test. Terraform prereleases are considered as potential candidates for
// the upcoming version and therefore are treated as semantically equal for
// testing purposes. If running prereleases of the same patch release is
// desired, give a different prerelease version. For example, if given
// version.Must(version.NewVersion("1.8.0-rc2")), then 1.8.0-rc1 will
// run the test.
func SkipIf(version *version.Version) TerraformVersionCheck {
	return skipIfCheck{
		version: version,
	}
}

// skipIfCheck implements the TerraformVersionCheck interface
type skipIfCheck struct {
	version *version.Version
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (s skipIfCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {
	var terraformVersion *version.Version

	// If given a prerelease version, check the Terraform CLI version directly,
	// otherwise use the core version so that prereleases are treated as equal.
	if s.version.Prerelease() != "" {
		terraformVersion = req.TerraformVersion
	} else {
		terraformVersion = req.TerraformVersion.Core()
	}

	if terraformVersion.Equal(s.version) {
		resp.Skip = fmt.Sprintf("Terraform CLI version is %s: skipping test.", s.version)
	}
}
