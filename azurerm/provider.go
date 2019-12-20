package azurerm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

func Provider() terraform.ResourceProvider {
	dataSources := map[string]*schema.Resource{
		"azurerm_client_config": dataSourceArmClientConfig(),
	}

	resources := map[string]*schema.Resource{

		"azurerm_autoscale_setting":   resourceArmAutoScaleSetting(),
		"azurerm_dashboard":           resourceArmDashboard(),
		"azurerm_management_lock":     resourceArmManagementLock(),
		"azurerm_template_deployment": resourceArmTemplateDeployment(),
	}

	return provider.AzureProvider(dataSources, resources)
}
