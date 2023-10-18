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

	if req.TerraformVersion.Equal(s.version) {
		resp.Skip = fmt.Sprintf("Terraform CLI version is %s: skipping test.", s.version)
	}
}
