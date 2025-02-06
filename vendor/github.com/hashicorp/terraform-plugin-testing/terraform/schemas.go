// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terraform

import (
	"github.com/hashicorp/terraform-plugin-testing/internal/configs/configschema"
)

// ProviderSchema represents the schema for a provider's own configuration
// and the configuration for some or all of its resources and data sources.
//
// The completeness of this structure depends on how it was constructed.
// When constructed for a configuration, it will generally include only
// resource types and data sources used by that configuration.
//
// Deprecated: This type is unintentionally exported by this Go module and not
// supported for external consumption. It will be removed in the next major
// version.
type ProviderSchema struct {
	Provider      *configschema.Block
	ResourceTypes map[string]*configschema.Block
	DataSources   map[string]*configschema.Block

	ResourceTypeSchemaVersions map[string]uint64
}

// ProviderSchemaRequest is used to describe to a ResourceProvider which
// aspects of schema are required, when calling the GetSchema method.
//
// Deprecated: This type is unintentionally exported by this Go module and not
// supported for external consumption. It will be removed in the next major
// version.
type ProviderSchemaRequest struct {
	ResourceTypes []string
	DataSources   []string
}
