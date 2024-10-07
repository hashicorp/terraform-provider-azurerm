// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"
	"errors"
	"strings"
)

// Any will return a nil error and empty skip message (run the test)
// if any of the given checks return a nil error and empty skip message.
// Otherwise, it will return all errors and fail the test if any of the given
// checks return a non-nil error, or it will return all skip messages
// and skip (pass) the test.
func Any(terraformVersionChecks ...TerraformVersionCheck) TerraformVersionCheck {
	return anyCheck{
		terraformVersionChecks: terraformVersionChecks,
	}
}

// anyCheck implements the TerraformVersionCheck interface
type anyCheck struct {
	terraformVersionChecks []TerraformVersionCheck
}

// CheckTerraformVersion satisfies the TerraformVersionCheck interface.
func (a anyCheck) CheckTerraformVersion(ctx context.Context, req CheckTerraformVersionRequest, resp *CheckTerraformVersionResponse) {
	var joinedErrors []error
	strBuilder := strings.Builder{}

	for _, subCheck := range a.terraformVersionChecks {
		checkResp := CheckTerraformVersionResponse{}

		subCheck.CheckTerraformVersion(ctx, CheckTerraformVersionRequest{TerraformVersion: req.TerraformVersion}, &checkResp)

		if checkResp.Error == nil && checkResp.Skip == "" {
			resp.Error = nil
			resp.Skip = ""
			return
		}

		joinedErrors = append(joinedErrors, checkResp.Error)

		if checkResp.Skip != "" {
			strBuilder.WriteString(checkResp.Skip)
			strBuilder.WriteString("\n")
		}
	}

	resp.Error = errors.Join(joinedErrors...)
	resp.Skip = strings.TrimSpace(strBuilder.String())
}
