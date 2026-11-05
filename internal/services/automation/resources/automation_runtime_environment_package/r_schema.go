// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_runtime_environment_package

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/packageresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"automation_runtime_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: packageresource.ValidateRuntimeEnvironmentID,
		},

		"content_uri": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"content_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9]+\.[0-9]+(\.[0-9]+){0,2}$`),
				"The `content_version` must be with 2 to 4 segments (e.g. `1.0`, `1.0.0`, or `1.0.0.0`).",
			),
		},

		"hash_algorithm": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"hash_value"},
		},

		"hash_value": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"hash_algorithm"},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"default": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"size_in_bytes": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
