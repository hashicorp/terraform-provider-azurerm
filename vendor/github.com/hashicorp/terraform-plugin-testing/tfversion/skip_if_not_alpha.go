// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"fmt"
	"strings"
)

// SkipIfNotAlpha will skip (pass) the test if the Terraform CLI
// version is not an alpha prerelease (for example, 1.10.0-alpha20241023).
//
// Alpha builds of Terraform include experimental features, so this version check
// can be used for acceptance testing of experimental features, such as deferred actions.
func SkipIfNotAlpha() TerraformVersionCheck {
	return skipIfNotAlphaCheck{}
}

// skipIfNotAlphaCheck implements the TerraformVersionCheck interface
type skipIfNotAlphaCheck struct{}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (s skipIfNotAlphaCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {
	if strings.Contains(req.TerraformVersion.Prerelease(), "alpha") {
		return
	}

	resp.Skip = fmt.Sprintf("Terraform CLI version %s is not an alpha build: skipping test.", req.TerraformVersion)
}
