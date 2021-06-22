package web

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceFunctionAppHostKeys() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceFunctionAppHostKeysRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"master_key": {
				Type:       pluginsdk.TypeString,
				Computed:   true,
				Sensitive:  true,
				Deprecated: "This property has been renamed to `primary_key` and will be removed in v3.0 of the provider in support of HashiCorp's inclusive language policy which can be found here: https://discuss.hashicorp.com/t/inclusive-language-changes",
			},

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_function_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"event_grid_extension_config_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceFunctionAppHostKeysRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	functionSettings, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(functionSettings.Response) {
			return fmt.Errorf("Error: AzureRM Function App %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Function App %q: %+v", name, err)
	}

	if functionSettings.ID == nil {
		return fmt.Errorf("cannot read ID for AzureRM Function App %q (Resource Group %q)", name, resourceGroup)
	}
	d.SetId(*functionSettings.ID)

	//lintignore:R006
	return pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		res, err := client.ListHostKeys(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return pluginsdk.NonRetryableError(fmt.Errorf("Error: AzureRM Function App %q (Resource Group %q) was not found", name, resourceGroup))
			}

			return pluginsdk.RetryableError(fmt.Errorf("Error making Read request on AzureRM Function App Hostkeys %q: %+v", name, err))
		}

		d.Set("master_key", res.MasterKey)
		d.Set("primary_key", res.MasterKey)

		defaultFunctionKey := ""
		if v, ok := res.FunctionKeys["default"]; ok {
			defaultFunctionKey = *v
		}
		d.Set("default_function_key", defaultFunctionKey)

		eventGridExtensionConfigKey := ""
		if v, ok := res.SystemKeys["eventgridextensionconfig_extension"]; ok {
			eventGridExtensionConfigKey = *v
		}
		d.Set("event_grid_extension_config_key", eventGridExtensionConfigKey)

		return nil
	})
}
