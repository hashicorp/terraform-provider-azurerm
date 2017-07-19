package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
)

func appServicePlanStateRefreshFunc(client *ArmClient, resourceGroupName string, appserviceplan string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.appServicePlansClient.Get(resourceGroupName, appserviceplan)
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in appServicePlanStateRefreshFunc to Azure ARM for App Service Plan '%s' (RG: '%s'): %s", appserviceplan, resourceGroupName, err)
		}

		return res, string(res.AppServicePlanProperties.ProvisioningState), nil
	}
}
