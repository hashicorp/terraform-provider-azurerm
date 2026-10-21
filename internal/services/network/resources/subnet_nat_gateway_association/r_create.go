// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet_nat_gateway_association

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetNatGatewayAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Subnets
	vnetClient := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subnetId, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return err
	}

	gatewayId, err := natgateways.ParseNatGatewayID(d.Get("nat_gateway_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(gatewayId.ID())
	defer locks.UnlockByID(gatewayId.ID())
	vnetId := commonids.NewVirtualNetworkID(subnetId.SubscriptionId, subnetId.ResourceGroupName, subnetId.VirtualNetworkName)
	locks.ByID(vnetId.ID())
	defer locks.UnlockByID(vnetId.ID())
	locks.ByID(subnetId.ID())
	defer locks.UnlockByID(subnetId.ID())

	subnet, err := client.Get(ctx, *subnetId, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(subnet.HttpResponse) {
			return fmt.Errorf("%s was not found", *subnetId)
		}
		return fmt.Errorf("retrieving %s: %+v", *subnetId, err)
	}

	if model := subnet.Model; model != nil {
		if props := model.Properties; props != nil {
			// check if the resources are imported
			if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if gateway := props.NatGateway; gateway != nil {
					if gateway.Id != nil && model.Id != nil {
						return tf.ImportAsExistsError("azurerm_subnet_nat_gateway_association", *model.Id)
					}
				}
			}
			props.NatGateway = &subnets.SubResource{
				Id: pointer.To(gatewayId.ID()),
			}
		}
	}

	// TODO: migrate this resource to a composite ID
	if err := client.CreateOrUpdateCallbackThenPoll(ctx, *subnetId, *subnet.Model, sdk.SetIDCallback(meta, subnetId, d)); err != nil {
		return fmt.Errorf("updating NAT Gateway Association for %s: %+v", *subnetId, err)
	}
	d.SetId(subnetId.ID())

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    helpers.SubnetProvisioningStateRefreshFunc(ctx, client, *subnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of subnet for NAT Gateway Association for %s: %+v", *subnetId, err)
	}

	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    helpers.VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for NAT Gateway Association for %s: %+v", *subnetId, err)
	}

	return resourceSubnetNatGatewayAssociationRead(d, meta)
}
