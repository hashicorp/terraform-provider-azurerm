package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPeerAsn() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPeerAsnRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.PeerAsnName(),
			},
		},
	}
}

func dataSourceArmPeerAsnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PeerAsnsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Peer Asn %q was not found", name)
		}

		return fmt.Errorf("failed to retrieve Peer Asn %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	return nil
}
