// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
)

// All will return the first non-nil error or non-empty skip message
// if any of the given checks return a non-nil error or non-empty skip message.
// Otherwise, it will return a nil error and empty skip message (run the test)
//
// Use of All is only necessary when used in conjunction with Any as the
// TerraformVersionChecks field automatically applies a logical AND.
func All(terraformVersionChecks ...TerraformVersionCheck) TerraformVersionCheck {
	return allCheck{
		terraformVersionChecks: terraformVersionChecks,
	}
}

// allCheck implements the TerraformVersionCheck interface
type allCheck struct {
	terraformVersionChecks []TerraformVersionCheck
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (a allCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {

	for _, subCheck := range a.terraformVersionChecks {
		checkResp := CheckTerraformVersionResponse{}

		subCheck.CheckTerraformVersion(ctx, CheckTerraformVersionRequest{TerraformVersion: req.TerraformVersion}, &checkResp)

		if checkResp.Error != nil {
			resp.Error = checkResp.Error
			return
		}

		if checkResp.Skip != "" {
			resp.Skip = checkResp.Skip
			return
		}
	}
}
