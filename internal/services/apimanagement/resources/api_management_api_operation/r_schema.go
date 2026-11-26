// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_operation

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementApiOperationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"operation_id": schemaz.SchemaApiManagementChildName("A unique identifier for this API Operation."),

		"api_name": schemaz.SchemaApiManagementApiName("The name of the API within the API Management Service where this API Operation should be created."),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"display_name": {
			Description: "The Display Name for this API Management Operation.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"method": {
			Description: "The HTTP Method used for this API Management Operation, like `GET`, `DELETE`, `PUT` or `POST` - but not limited to these values.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"url_template": {
			Description: "The relative URL Template identifying the target resource for this operation, which may include parameters.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"description": {
			Description: "A description for this API Operation, which may include HTML formatting tags.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"request": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Description: "A description of the HTTP Request, which may include HTML tags.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
						Description: "The HTTP Status Code.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"description": {
						Description: "A description of the HTTP Response, which may include HTML tags.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"header": schemaz.SchemaApiManagementOperationParameterContract(),

					"representation": schemaz.SchemaApiManagementOperationRepresentation(),
				},
			},
		},

		"template_parameter": schemaz.SchemaApiManagementOperationParameterContract(),
	}
}
