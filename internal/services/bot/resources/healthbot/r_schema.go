// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package healthbot

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthbot/2025-05-25/healthbots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func healthbotSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Healthbot Service resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.HealthbotName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Description:  "The name which should be used for the SKU of the service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(healthbots.PossibleValuesForSkuName(), false),
		},

		"bot_management_portal_url": {
			Description: "The management portal URL.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
