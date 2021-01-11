package automation

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceAutomationAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAutomationAccountRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"primary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAutomationAccountRead(d *schema.ResourceData, meta interface{}) error {
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
