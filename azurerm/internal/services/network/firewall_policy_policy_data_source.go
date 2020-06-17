package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmFirewallPolicyPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmFirewallPolicyPolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.FirewallPolicyPolicyName(),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmFirewallPolicyPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Firewall Policy Policy %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("retrieving Firewall Policy Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Firewall Policy Policy %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	return tags.FlattenAndSet(d, resp.Tags)
}
