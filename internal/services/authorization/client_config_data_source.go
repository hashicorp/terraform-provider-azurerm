package authorization

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceArmClientConfig() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmClientConfigRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"client_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"subscription_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"object_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmClientConfigRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client)
	_, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	d.SetId(time.Now().UTC().String())
	d.Set("client_id", client.Account.ClientId)
	d.Set("object_id", client.Account.ObjectId)
	d.Set("subscription_id", client.Account.SubscriptionId)
	d.Set("tenant_id", client.Account.TenantId)

	return nil
}
