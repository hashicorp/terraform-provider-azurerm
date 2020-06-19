package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmFirewallPolicyRuleCollectionGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmFirewallPolicyRuleCollectionGroupRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.FirewallPolicyRuleCollectionGroupName(),
			},

			"firewall_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FirewallPolicyID,
			},
		},
	}
}

func dataSourceArmFirewallPolicyRuleCollectionGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPolicyRuleGroupClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	policyId, err := parse.FirewallPolicyID(d.Get("firewall_policy_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, policyId.ResourceGroup, policyId.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Firewall Policy Rule Collection Group %q (Resource Group %q / Policy %q) was not found", name, policyId.ResourceGroup, policyId.Name)
		}

		return fmt.Errorf("retrieving Firewall Policy Rule Collection Group %q (Resource Group %q / Policy %q): %+v", name, policyId.ResourceGroup, policyId.Name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Firewall Policy Rule Collection Group %q (Resource Group %q / Policy %q) ID", name, policyId.ResourceGroup, policyId.Name)
	}

	d.SetId(*resp.ID)

	return nil
}
