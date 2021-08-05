package validate

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	managementGroupValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"
	resourceValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	subscriptionValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func ResourceAssignmentId() pluginsdk.SchemaValidateFunc {
	return validation.All(
		validation.None(
			map[string]func(interface{}, string) ([]string, []error){
				"Management Group ID": managementGroupValidate.ManagementGroupID,
				"Resource Group ID":   resourceValidate.ResourceGroupID,
				"Subscription ID":     subscriptionValidate.SubscriptionID,
			},
		),
		azure.ValidateResourceID,
	)
}
