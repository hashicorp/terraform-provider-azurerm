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
// For example, if given a minimum version of version.Must(version.NewVersion("0.15.0"))
// and a maximum version of version.Must(version.NewVersion("0.16.0")), then versions 0.15.x
// will skip the test.
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

	if req.TerraformVersion.GreaterThanOrEqual(s.minimumVersion) && req.TerraformVersion.LessThan(s.maximumVersion) {
		resp.Skip = fmt.Sprintf("Terraform CLI version %s is between %s and %s: skipping test.",
			req.TerraformVersion, s.minimumVersion, s.maximumVersion)
	}
}
