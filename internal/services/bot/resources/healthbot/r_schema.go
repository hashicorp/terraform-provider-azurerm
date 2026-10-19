// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.HealthbotName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(healthbots.PossibleValuesForSkuName(), false),
		},

		"bot_management_portal_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
