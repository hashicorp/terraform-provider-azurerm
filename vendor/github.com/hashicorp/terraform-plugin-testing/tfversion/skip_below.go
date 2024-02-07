// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-version"
)

// SkipBelow will skip (pass) the test if the Terraform CLI
// version is below the given version. For example, if given
// version.Must(version.NewVersion("0.15.0")), then 0.14.x or
// any other prior minor versions will skip the test.
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

	if req.TerraformVersion.LessThan(s.minimumVersion) {
		resp.Skip = fmt.Sprintf("Terraform CLI version %s is below minimum version %s: skipping test",
			req.TerraformVersion, s.minimumVersion)
	}
}
