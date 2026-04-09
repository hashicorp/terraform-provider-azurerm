// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import (
	"context"

	"github.com/hashicorp/go-version"
)

// TerraformVersionCheck is the interface for writing check logic against the Terraform CLI version.
// The Terraform CLI version is determined by the binary selected by the TF_ACC_TERRAFORM_PATH environment
// variable value, installed by the TF_ACC_TERRAFORM_VERSION value, or already existing based on the PATH environment
// variable. This logic is executed at the beginning of the TestCase before any TestStep is executed.
//
// This package contains some built-in functionality that implements the interface, otherwise consumers can use this
// interface for implementing their own custom logic.
type TerraformVersionCheck interface {
	// CheckTerraformVersion should implement the logic to either pass, error (failing the test), or skip (passing the test).
	CheckTerraformVersion(context.Context, CheckTerraformVersionRequest, *CheckTerraformVersionResponse)
}

// CheckTerraformVersionRequest is the request received for the CheckTerraformVersion method of the
// TerraformVersionCheck interface. The response of that method is CheckTerraformVersionResponse.
type CheckTerraformVersionRequest struct {
	// TerraformVersion is the version associated with the selected Terraform CLI binary.
	TerraformVersion *version.Version
}

// CheckTerraformVersionResponse is the response returned for the CheckTerraformVersion method of the
// TerraformVersionCheck interface. The request of that method is CheckTerraformVersionRequest.
type CheckTerraformVersionResponse struct {
	// Error will result in failing the test with a given error message.
	Error error

	// Skip will result in passing the test with a given skip message.
	Skip string
}
