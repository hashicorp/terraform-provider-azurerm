// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"path/filepath"
	"strconv"
)

// StaticFile returns the supplied file.
func StaticFile(file string) func(TestStepConfigRequest) string {
	return func(_ TestStepConfigRequest) string {
		return file
	}
}

// TestNameFile returns the name of the test suffixed with the supplied
// file and prefixed with "testdata".
//
// For example, given test code:
//
//	func TestExampleCloudThing_basic(t *testing.T) {
//	    resource.Test(t, resource.TestCase{
//	        Steps: []resource.TestStep{
//	            {
//	                ConfigFile: config.TestNameFile("test.tf"),
//	            },
//	        },
//	    })
//	}
//
// The testing configuration will be expected in the
// testdata/TestExampleCloudThing_basic/test.tf file.
func TestNameFile(file string) func(TestStepConfigRequest) string {
	return func(req TestStepConfigRequest) string {
		return filepath.Join("testdata", req.TestName, file)
	}
}

// TestStepFile returns the name of the test suffixed with the test
// step number and the supplied file, and prefixed with "testdata".
//
// For example, given test code:
//
//	func TestExampleCloudThing_basic(t *testing.T) {
//	    resource.Test(t, resource.TestCase{
//	        Steps: []resource.TestStep{
//	            {
//	                ConfigFile: config.TestStepFile("test.tf"),
//	            },
//	        },
//	    })
//	}
//
// The testing configuration will be expected in the
// testdata/TestExampleCloudThing_basic/1/test.tf file
// as TestStepConfigRequest.StepNumber is one-based.
func TestStepFile(file string) func(TestStepConfigRequest) string {
	return func(req TestStepConfigRequest) string {
		return filepath.Join("testdata", req.TestName, strconv.Itoa(req.StepNumber), file)
	}
}
