// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terraform

// ResourceType is a type of resource that a resource provider can manage.
//
// Deprecated: This type is unintentionally exported by this Go module and not
// supported for external consumption. It will be removed in the next major
// version.
type ResourceType struct {
	Name       string // Name of the resource, example "instance" (no provider prefix)
	Importable bool   // Whether this resource supports importing

	// SchemaAvailable is set if the provider supports the ProviderSchema,
	// ResourceTypeSchema and DataSourceSchema methods. Although it is
	// included on each resource type, it's actually a provider-wide setting
	// that's smuggled here only because that avoids a breaking change to
	// the plugin protocol.
	SchemaAvailable bool
}

// DataSource is a data source that a resource provider implements.
//
// Deprecated: This type is unintentionally exported by this Go module and not
// supported for external consumption. It will be removed in the next major
// version.
type DataSource struct {
	Name string

	// SchemaAvailable is set if the provider supports the ProviderSchema,
	// ResourceTypeSchema and DataSourceSchema methods. Although it is
	// included on each resource type, it's actually a provider-wide setting
	// that's smuggled here only because that avoids a breaking change to
	// the plugin protocol.
	SchemaAvailable bool
}
