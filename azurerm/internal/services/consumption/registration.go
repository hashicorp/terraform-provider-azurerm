package consumption

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// The Consumption Budget resource names are extracted into their own variables
	// as the core logic for the Consumption Budget resources is generic and has been
	// extracted out of the specific Consumption Budget resources. These constants are
	// used when the generic Consumption Budget functions require a resource name.
	consumptionBudgetResourceGroupName = "azurerm_consumption_budget_resource_group"
	consumptionBudgetSubscriptionName  = "azurerm_consumption_budget_subscription"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Consumption"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Consumption",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		consumptionBudgetResourceGroupName: resourceArmConsumptionBudgetResourceGroup(),
		consumptionBudgetSubscriptionName:  resourceArmConsumptionBudgetSubscription(),
	}
}
