package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	managementGroupValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	resourceValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	subscriptionValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
