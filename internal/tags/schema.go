// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

import "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"

// SchemaDataSource returns the Schema which should be used for Tags on a Data Source
func SchemaDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeMap,
		Computed: true,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

// ForceNewSchema returns the Schema which should be used for Tags when changes
// require recreation of the resource
func ForceNewSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeMap,
		Optional:     true,
		ForceNew:     true,
		ValidateFunc: Validate,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

// Schema returns the Schema used for Tags
func Schema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeMap,
		Optional:     true,
		ValidateFunc: Validate,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

// SchemaWithMax returns the Schema with the maximum used for Tags
func SchemaWithMax(max int) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeMap,
		Optional:     true,
		ValidateFunc: ValidateWithMax(max),
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

// Schema returns the Schema used for Tags
func SchemaEnforceLowerCaseKeys() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeMap,
		Optional:     true,
		ValidateFunc: EnforceLowerCaseKeys,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}
