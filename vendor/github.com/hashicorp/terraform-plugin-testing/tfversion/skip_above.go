// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// SkipAbove will skip (pass) the test if the Terraform CLI
// version is below the given version. For example, if given
// version.Must(version.NewVersion("0.15.0")), then 0.14.x or
// any other prior minor versions will skip the test.
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

	if req.TerraformVersion.GreaterThan(s.maximumVersion) {
		resp.Skip = fmt.Sprintf("Terraform CLI version %s is above maximum version %s: skipping test",
			req.TerraformVersion, s.maximumVersion)
	}
}
