// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfversion

import "github.com/hashicorp/go-version"

// Common use version variables to simplify provider testing implementations.
// This list is not intended to be exhaustive of all Terraform versions,
// however these should at least include cases where Terraform
// introduced new configuration language features.
var (
	// Version0_12_26 is the first Terraform CLI version supported
	// by the testing code.
	Version0_12_26 *version.Version = version.Must(version.NewVersion("0.12.26"))

	// Major versions

	Version1_0_0 *version.Version = version.Must(version.NewVersion("1.0.0"))
	Version2_0_0 *version.Version = version.Must(version.NewVersion("2.0.0"))

	// Minor versions

	Version0_13_0 *version.Version = version.Must(version.NewVersion("0.13.0"))
	Version0_14_0 *version.Version = version.Must(version.NewVersion("0.14.0"))
	Version0_15_0 *version.Version = version.Must(version.NewVersion("0.15.0"))
	Version1_1_0  *version.Version = version.Must(version.NewVersion("1.1.0"))
	Version1_2_0  *version.Version = version.Must(version.NewVersion("1.2.0"))
	Version1_3_0  *version.Version = version.Must(version.NewVersion("1.3.0"))
	Version1_4_0  *version.Version = version.Must(version.NewVersion("1.4.0"))
	Version1_5_0  *version.Version = version.Must(version.NewVersion("1.5.0"))
)
