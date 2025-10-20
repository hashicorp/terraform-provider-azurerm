// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

// TestStepConfigFunc is the callback type used with acceptance tests to
// specify a string which either identifies a directory containing
// Terraform configuration files, or a file that contains Terraform
// configuration.
type TestStepConfigFunc func(TestStepConfigRequest) string

// TestStepConfigRequest defines the request supplied to types
// implementing TestStepConfigFunc. StepNumber is one-based
// and is used in the predefined helper functions:
//
//   - [config.TestStepDirectory]
//   - [config.TestStepFile].
//
// TestName is used in the predefined helper functions:
//
//   - [config.TestNameDirectory]
//   - [config.TestStepDirectory]
//   - [config.TestNameFile]
//   - [config.TestStepFile]
type TestStepConfigRequest struct {
	StepNumber int
	TestName   string
}

// Exec executes TestStepConfigFunc if it is not nil, otherwise an
// empty string is returned.
func (f TestStepConfigFunc) Exec(req TestStepConfigRequest) string {
	if f != nil {
		return f(req)
	}

	return ""
}
