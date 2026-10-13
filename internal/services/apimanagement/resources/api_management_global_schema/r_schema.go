// Copyright IBM Corp. 2014, 2025
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
		"schema_id": schemaz.SchemaApiManagementChildName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(schema.PossibleValuesForSchemaType(), false),
		},

		"value": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringIsNotEmpty,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
