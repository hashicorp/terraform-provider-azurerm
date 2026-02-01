// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov6

// ResourceIdentityData contains the raw undecoded identity data
// for a resource.
type ResourceIdentityData struct {
	// IdentityData is represented as a `DynamicValue`. See the documentation for
	// `DynamicValue` for information about safely creating the
	// `DynamicValue`.
	// The identity should be represented as a tftypes.Object, with each
	// attribute and nested block getting its own key and value.
	IdentityData *DynamicValue
}
