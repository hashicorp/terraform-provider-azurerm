package automation

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceAutomationAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAutomationAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"primary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"secondary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAutomationAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	iclient := meta.(*clients.Client).Automation.AgentRegistrationInfoClient
	client := meta.(*clients.Client).Automation.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Automation Account %q (Resource Group %q) was not found", name, resourceGroupName)
		}
		return fmt.Errorf("Error making Read request on Automation %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}
	d.SetId(*resp.ID)

	iresp, err := iclient.Get(ctx, resourceGroupName, name)
	if err != nil {
		if utils.ResponseWasNotFound(iresp.Response) {
			return fmt.Errorf("Error: Automation Account Registration Information %q (Resource Group %q) was not found", name, resourceGroupName)
		}
		return fmt.Errorf("Error making Read request on Automation Account Registration Information %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}
	d.Set("primary_key", iresp.Keys.Primary)
	d.Set("secondary_key", iresp.Keys.Secondary)
	d.Set("endpoint", iresp.Endpoint)
	return nil
}
