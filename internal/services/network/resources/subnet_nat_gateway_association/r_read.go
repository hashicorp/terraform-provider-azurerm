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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetNatGatewayAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Subnets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	subnet, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(subnet.HttpResponse) {
			log.Printf("[DEBUG] %s could not be found - removing from state!", *id)
			d.SetId("")
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
		log.Printf("[DEBUG] %s doesn't have a NAT Gateway - removing from state!", *id)
		d.SetId("")
		return nil
	}

	gatewayId, err := natgateways.ParseNatGatewayID(*props.NatGateway.Id)
	if err != nil {
		return err
	}

	d.Set("subnet_id", id.ID())
	d.Set("nat_gateway_id", gatewayId.ID())

	return nil
}
