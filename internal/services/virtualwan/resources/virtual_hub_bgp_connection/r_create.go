// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_bgp_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubBgpConnectionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	virtHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtHubId.VirtualHubName, virtualHubResourceName)

	id := commonids.NewVirtualHubBGPConnectionID(virtHubId.SubscriptionId, virtHubId.ResourceGroupName, virtHubId.VirtualHubName, d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.VirtualHubBgpConnectionGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_hub_bgp_connection", id.ID())
		}
	}

	parameters := virtualwans.BgpConnection{
		Name: pointer.To(d.Get("name").(string)),
		Properties: &virtualwans.BgpConnectionProperties{
			PeerAsn: pointer.To(int64(d.Get("peer_asn").(int))),
			PeerIP:  pointer.To(d.Get("peer_ip").(string)),
		},
	}

	if v, ok := d.GetOk("virtual_network_connection_id"); ok {
		parameters.Properties.HubVirtualNetworkConnection = &virtualwans.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if err := client.VirtualHubBgpConnectionCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceVirtualHubBgpConnectionRead(d, meta)
}
