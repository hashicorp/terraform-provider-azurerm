// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package point_to_site_vpn_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePointToSiteVPNGatewayDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualWANP2SVPNGatewayID(d.Id())
	if err != nil {
		return err
	}

	if err := client.P2sVpnGatewaysDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
