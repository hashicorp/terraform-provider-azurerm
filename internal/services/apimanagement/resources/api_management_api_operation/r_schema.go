// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementApiOperationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"operation_id": schemaz.SchemaApiManagementChildName(),

		"api_name": schemaz.SchemaApiManagementApiName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"method": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"url_template": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"request": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"header": schemaz.SchemaApiManagementOperationParameterContract(),

					"query_parameter": schemaz.SchemaApiManagementOperationParameterContract(),

					"representation": schemaz.SchemaApiManagementOperationRepresentation(),
				},
			},
		},

		"response": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"status_code": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"description": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"header": schemaz.SchemaApiManagementOperationParameterContract(),

					"representation": schemaz.SchemaApiManagementOperationRepresentation(),
				},
			},
		},

		"template_parameter": schemaz.SchemaApiManagementOperationParameterContract(),
	}
}
