// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"
)

// SkipIfNotPrerelease will skip (pass) the test if the Terraform CLI
// version does not include prerelease information. This will include builds
// of Terraform that are from source. (e.g. 1.8.0-dev)
func SkipIfNotPrerelease() TerraformVersionCheck {
	return skipIfNotPrereleaseCheck{}
}

// skipIfNotPrereleaseCheck implements the TerraformVersionCheck interface
type skipIfNotPrereleaseCheck struct{}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (s skipIfNotPrereleaseCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {
	if req.TerraformVersion.Prerelease() != "" {
		return
	}

	resp.Skip = fmt.Sprintf("Terraform CLI version %s is not a prerelease build: skipping test.", req.TerraformVersion)
}
