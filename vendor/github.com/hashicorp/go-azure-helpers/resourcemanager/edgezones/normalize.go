// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package edgezones

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
)

// Normalize transforms the specified user input into a canonical form
func Normalize(input string) string {
	// we're intentionally passing through to Locations today since this is sufficient
	// but it's helpful to have a specific endpoint for this should this need to change
	// in the future
	return location.Normalize(input)
}

// NormalizeNilable normalizes the specified user input into a canonical form
func NormalizeNilable(input *string) string {
	// we're intentionally passing through to Locations today since this is sufficient
	// but it's helpful to have a specific endpoint for this should this need to change
	// in the future
	return location.NormalizeNilable(input)
}
