package costmanagement

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SubscriptionCostManagementScheduledActionResource struct {
	base costManagementScheduledActionBaseResource
}

var _ sdk.Resource = SubscriptionCostManagementScheduledActionResource{}

func (r SubscriptionCostManagementScheduledActionResource) Arguments() map[string]*pluginsdk.Schema {
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

func (r SubscriptionCostManagementScheduledActionResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r SubscriptionCostManagementScheduledActionResource) ModelObject() interface{} {
	return nil
}

func (r SubscriptionCostManagementScheduledActionResource) ResourceType() string {
	return "azurerm_subscription_cost_management_scheduled_action"
}

func (r SubscriptionCostManagementScheduledActionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SubscriptionCostManagementScheduledActionID
}

func (r SubscriptionCostManagementScheduledActionResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "subscription_id")
}

func (r SubscriptionCostManagementScheduledActionResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("subscription_id")
}

func (r SubscriptionCostManagementScheduledActionResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r SubscriptionCostManagementScheduledActionResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}
