// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_schema

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementApiSchemaSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"schema_id": schemaz.SchemaApiManagementChildName(),

		"api_name": schemaz.SchemaApiManagementApiName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"content_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"value": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				if d.Get("content_type") == "application/vnd.ms-azure-apim.swagger.definitions+json" || d.Get("content_type") == "application/vnd.oai.openapi.components+json" {
					return pluginsdk.SuppressJsonDiff(k, old, new, d)
				}
				return old == new
			},
			ExactlyOneOf: []string{"value", "definitions", "components"},
		},

		"components": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			ExactlyOneOf:     []string{"value", "definitions", "components"},
		},

		"definitions": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			ExactlyOneOf:     []string{"value", "definitions", "components"},
		},
	}
}
