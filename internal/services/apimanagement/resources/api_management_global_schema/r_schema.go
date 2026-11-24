// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_global_schema

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementGlobalSchemaSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"schema_id": schemaz.SchemaApiManagementChildName("A unique identifier for this Schema."),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"type": {
			Description:  "The content type of the Schema.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(schema.PossibleValuesForSchemaType(), false),
		},

		"value": {
			Description:      "The string defining the document representing the Schema.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringIsNotEmpty,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"description": {
			Description:  "The description of the schema.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
