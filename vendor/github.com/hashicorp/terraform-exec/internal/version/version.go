// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package version

const version = "0.25.1"

// ModuleVersion returns the current version of the github.com/hashicorp/terraform-exec Go module.
// This is a function to allow for future possible enhancement using debug.BuildInfo.
func ModuleVersion() string {
	return version
}
