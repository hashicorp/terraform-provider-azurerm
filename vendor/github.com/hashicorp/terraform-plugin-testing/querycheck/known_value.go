// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package querycheck

import (
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// KnownValueCheck represents a check of a known value at a specific JSON path
// and is used to specify multiple known value checks to assert against a
// single resource object returned by a query.
type KnownValueCheck struct {
	// Path specifies the JSON path to check within the resource object.
	Path tfjsonpath.Path
	// KnownValue specifies the expected known value check to perform at the given path.
	KnownValue knownvalue.Check
}
