// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet_nat_gateway_association

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetNatGatewayAssociationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Subnets
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	subnet, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(subnet.HttpResponse) {
			log.Printf("[DEBUG] %s could not be found - removing from state!", *id)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if subnet.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil ", id)
	}
	if subnet.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil ", id)
	}

	props := subnet.Model.Properties

	if props.NatGateway == nil || props.NatGateway.Id == nil {
		log.Printf("[DEBUG] %s has no NAT Gateway - removing from state!", *id)
		return nil
	}
	gatewayId, err := natgateways.ParseNatGatewayID(*props.NatGateway.Id)
	if err != nil {
		return err
	}

	locks.ByID(gatewayId.ID())
	defer locks.UnlockByID(gatewayId.ID())
	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	locks.ByID(vnetId.ID())
	defer locks.UnlockByID(vnetId.ID())

	subnet, err = client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(subnet.HttpResponse) {
			log.Printf("[DEBUG] %s could not be found - removing from state!", *id)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	subnet.Model.Properties.NatGateway = nil

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *subnet.Model); err != nil {
		return fmt.Errorf("removing %s: %+v", *id, err)
	}

	return nil
}
