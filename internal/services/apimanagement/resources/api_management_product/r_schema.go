// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_product

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementProductSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"product_id": schemaz.SchemaApiManagementChildName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"subscription_required": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"published": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"approval_required": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"terms": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"subscriptions_limit": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},
	}
}
