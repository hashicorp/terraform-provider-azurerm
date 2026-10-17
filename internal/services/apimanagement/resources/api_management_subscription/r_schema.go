// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_subscription

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/subscription"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementSubscriptionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.Any(validate.ApiManagementChildName, validation.StringIsEmpty),
		},

		// 3.0 this seems to have been renamed to owner id?
		"user_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_management_name": schemaz.SchemaApiManagementName(),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"product_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  product.ValidateProductID,
			ConflictsWith: []string{"api_id"},
		},

		"api_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  api.ValidateApiID,
			ConflictsWith: []string{"product_id"},
		},

		"state": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(subscription.SubscriptionStateSubmitted),
			ValidateFunc: validation.StringInSlice(subscription.PossibleValuesForSubscriptionState(), false),
		},

		"primary_key": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because the API generates a subscription key when not specified
			Computed:  true,
			Sensitive: true,
		},

		"secondary_key": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because the API generates a subscription key when not specified
			Computed:  true,
			Sensitive: true,
		},

		"allow_tracing": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
	}
}
