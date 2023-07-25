// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SubscriptionCostManagementViewResource struct {
	base costManagementViewBaseResource
}

var _ sdk.Resource = SubscriptionCostManagementViewResource{}

func (r SubscriptionCostManagementViewResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},
	}
	return r.base.arguments(schema)
}

func (r SubscriptionCostManagementViewResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r SubscriptionCostManagementViewResource) ModelObject() interface{} {
	return nil
}

func (r SubscriptionCostManagementViewResource) ResourceType() string {
	return "azurerm_subscription_cost_management_view"
}

func (r SubscriptionCostManagementViewResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SubscriptionCostManagementViewID
}

func (r SubscriptionCostManagementViewResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "subscription_id")
}

func (r SubscriptionCostManagementViewResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("subscription_id")
}

func (r SubscriptionCostManagementViewResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r SubscriptionCostManagementViewResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}
