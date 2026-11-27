// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_connection

import (
	"fmt"
	"time"

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

func resourceVirtualHubConnectionCreateOrUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	virtualHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	id := virtualwans.NewHubVirtualNetworkConnectionID(virtualHubId.SubscriptionId, virtualHubId.ResourceGroupName, virtualHubId.VirtualHubName, d.Get("name").(string))

	locks.ByName(virtualHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubId.VirtualHubName, virtualHubResourceName)

	remoteVirtualNetworkId, err := commonids.ParseVirtualNetworkID(d.Get("remote_virtual_network_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(remoteVirtualNetworkId.ID())
	defer locks.UnlockByID(remoteVirtualNetworkId.ID())

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.HubVirtualNetworkConnectionsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(azureVirtualHubConnectionResourceName, id.ID())
			}
		}
	}

	connection := virtualwans.HubVirtualNetworkConnection{
		Name: pointer.To(id.HubVirtualNetworkConnectionName),
		Properties: &virtualwans.HubVirtualNetworkConnectionProperties{
			RemoteVirtualNetwork: &virtualwans.SubResource{
				Id: pointer.To(remoteVirtualNetworkId.ID()),
			},
			EnableInternetSecurity: pointer.To(d.Get("internet_security_enabled").(bool)),
		},
	}

	if v, ok := d.GetOk("routing"); ok {
		connection.Properties.RoutingConfiguration = expandVirtualHubConnectionRouting(v.([]any))
	}

	if err := client.HubVirtualNetworkConnectionsCreateOrUpdateCallbackThenPoll(ctx, id, connection, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	timeout, _ := ctx.Deadline()

	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(virtualwans.ProvisioningStateUpdating)},
		Target:     []string{string(virtualwans.ProvisioningStateSucceeded)},
		Refresh:    virtualHubConnectionProvisioningStateRefreshFunc(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	return resourceVirtualHubConnectionRead(d, meta)
}
