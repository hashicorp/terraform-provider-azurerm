// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_bgp_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubBgpConnectionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := commonids.ParseVirtualHubBGPConnectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.HubName, virtualHubResourceName)
	defer locks.UnlockByName(id.HubName, virtualHubResourceName)

	if err := client.VirtualHubBgpConnectionDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
