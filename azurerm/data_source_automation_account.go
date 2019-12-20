package azurerm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmAutomationAccount() *schema.Resource {
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

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

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
	client := meta.(*clients.Client).Automation.AgentRegistrationInfoClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Automation Account %q (Resource Group %q) was not found", name, resourceGroupName)
		}
		return fmt.Errorf("Error making Read request on Automation Account Registration Information %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}
	d.SetId(*resp.ID)
	d.Set("primary_key", resp.Keys.Primary)
	d.Set("secondary_key", resp.Keys.Secondary)
	d.Set("endpoint", resp.Endpoint)
	return nil
}
