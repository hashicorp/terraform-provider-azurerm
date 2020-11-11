package subscription

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSubscriptionAlias() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSubscriptionAliasRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"subscription_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmSubscriptionAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Subscription.AliasClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Subscription Alias %q does not exist", name)
		}
		return fmt.Errorf("retrieving Subscription Alias %q : %+v", name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Subscription Alias %q  ID", name)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	if prop := resp.Properties; prop != nil {
		d.Set("subscription_id", prop.SubscriptionID)
	}
	return nil
}
