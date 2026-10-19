// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package consumption_budget_subscription

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ConsumptionBudgetName(),
		},
		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},
	}
	return r.base.Arguments(schema)
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.Attributes()
}
