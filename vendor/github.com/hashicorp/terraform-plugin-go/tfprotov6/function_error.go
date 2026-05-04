// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov6

// FunctionError is used to convey information back to the user running Terraform.
type FunctionError struct {
	// Text is the description of the error.
	Text string

	// FunctionArgument is the positional function argument for aligning
	// configuration source.
	FunctionArgument *int64
}
