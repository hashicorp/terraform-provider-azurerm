// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azure

import "strings"

// @tombuildsstuff: this should be somewhere else (probably go-azure-helpers) but this'll do temporarily

func TitleCase(input string) string {
	//lint:ignore SA1019
	return strings.Title(input) //nolint:staticcheck
}
