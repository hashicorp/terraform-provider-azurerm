// Copyright IBM Corp. 2014, 2026
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
		"product_id": schemaz.SchemaApiManagementChildName("The Identifier for this Product, which must be unique within the API Management Service."),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"display_name": {
			Description:  "The Display Name for this API Management Product.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"subscription_required": {
			Description: "Is a Subscription required to access API's included in this Product?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"published": {
			Description: "Is this Product Published?",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"approval_required": {
			Description: "Do subscribers need to be approved prior to being able to use the Product?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"description": {
			Description: "A description of this Product, which may include HTML formatting tags.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"terms": {
			Description: "The Terms and Conditions for this Product, which must be accepted by Developers before they can begin the Subscription process.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"subscriptions_limit": {
			Description: "The number of subscriptions a user can have to this Product at the same time.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
		},
	}
}
