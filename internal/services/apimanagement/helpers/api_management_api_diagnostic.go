// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apidiagnostic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ResourceApiManagementApiDiagnosticAdditionalContentSchema() *pluginsdk.Schema {
	// lintignore:XS003
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Computed: true, // azignore:AZS007 - pre-existing violation
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"body_bytes": {
					Description:  "Number of payload bytes to log (up to 8192).",
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 8192),
				},
				"headers_to_log": {
					Description: "Specifies a list of headers to log.",
					Type:        pluginsdk.TypeSet,
					Optional:    true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Set: pluginsdk.HashString,
				},
				"data_masking": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"query_params": SchemaApiManagementDataMaskingEntityList(),
							"headers":      SchemaApiManagementDataMaskingEntityList(),
						},
					},
				},
			},
		},
	}
}

func SchemaApiManagementDataMaskingEntityList() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"mode": {
					Description:  "The data masking mode.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(apidiagnostic.PossibleValuesForDataMaskingMode(), false),
				},

				"value": {
					Description:  "The name of the header or the query parameter to mask.",
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}
