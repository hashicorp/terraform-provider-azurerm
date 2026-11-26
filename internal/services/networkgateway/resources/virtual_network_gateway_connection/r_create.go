// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_connection

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgatewayconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkGatewayConnectionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetworkGateway.VirtualNetworkGatewayConnections
	vnetGatewayClient := meta.(*clients.Client).NetworkGateway.VirtualNetworkGateways
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualnetworkgatewayconnections.NewConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_network_gateway_connection", id.ID())
		}
	}

	var virtualNetworkGateway virtualnetworkgateways.VirtualNetworkGateway
	if v, ok := d.GetOk("virtual_network_gateway_id"); ok {
		virtualNetworkGatewayId := v.(string)

		gwid, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(virtualNetworkGatewayId)
		if err != nil {
			return err
		}

		resp, err := vnetGatewayClient.Get(ctx, *gwid)
		if err != nil {
			return err
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: %+v", gwid, err)
		}

		virtualNetworkGateway = *resp.Model
	}

	properties, err := getVirtualNetworkGatewayConnectionProperties(d, virtualNetworkGateway)
	if err != nil {
		return err
	}

	connection := virtualnetworkgatewayconnections.VirtualNetworkGatewayConnection{
		Name:       pointer.To(id.ConnectionName),
		Location:   pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
		Properties: *properties,
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, connection, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	if properties.SharedKey != nil && !d.IsNewResource() {
		if err := client.SetSharedKeyThenPoll(ctx, id, virtualnetworkgatewayconnections.ConnectionSharedKey{
			Value: pointer.From(properties.SharedKey),
		}); err != nil {
			return fmt.Errorf("updating Shared Key for %s: %+v", id, err)
		}

		// Once this issue https://github.com/Azure/azure-rest-api-specs/issues/26660 is fixed, below this part will be removed
		stateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{string(virtualnetworkgatewayconnections.ProvisioningStateUpdating)},
			Target:     []string{string(virtualnetworkgatewayconnections.ProvisioningStateSucceeded)},
			Refresh:    virtualNetworkGatewayConnectionStateRefreshFunc(ctx, client, id),
			MinTimeout: 15 * time.Second,
			Timeout:    d.Timeout(pluginsdk.TimeoutUpdate),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for update of %s: %+v", id, err)
		}
	}

	return resourceVirtualNetworkGatewayConnectionRead(d, meta)
}
