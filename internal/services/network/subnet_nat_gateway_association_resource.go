// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/natgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetNatGatewayAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSubnetNatGatewayAssociationCreate,
		Read:   resourceSubnetNatGatewayAssociationRead,
		Delete: resourceSubnetNatGatewayAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseSubnetID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"nat_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: natgateways.ValidateNatGatewayID,
			},
		},
	}
}

func resourceSubnetNatGatewayAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.Subnets
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

	locks.ByName(gatewayId.NatGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(gatewayId.NatGatewayName, natGatewayResourceName)
	locks.ByName(subnetId.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(subnetId.VirtualNetworkName, VirtualNetworkResourceName)
	locks.ByName(subnetId.SubnetName, SubnetResourceName)
	defer locks.UnlockByName(subnetId.SubnetName, SubnetResourceName)

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
			if gateway := props.NatGateway; gateway != nil {
				if gateway.Id != nil && model.Id != nil {
					return tf.ImportAsExistsError("azurerm_subnet_nat_gateway_association", *model.Id)
				}
			}
			props.NatGateway = &subnets.SubResource{
				Id: pointer.To(gatewayId.ID()),
			}
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *subnetId, *subnet.Model); err != nil {
		return fmt.Errorf("updating NAT Gateway Association for %s: %+v", *subnetId, err)
	}

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    SubnetProvisioningStateRefreshFunc(ctx, client, *subnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of subnet for NAT Gateway Association for %s: %+v", *subnetId, err)
	}

	vnetId := commonids.NewVirtualNetworkID(subnetId.SubscriptionId, subnetId.ResourceGroupName, subnetId.VirtualNetworkName)
	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for NAT Gateway Association for %s: %+v", *subnetId, err)
	}

	d.SetId(subnetId.ID())

	return resourceSubnetNatGatewayAssociationRead(d, meta)
}

func resourceSubnetNatGatewayAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.Subnets
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

func resourceSubnetNatGatewayAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.Subnets
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

	locks.ByName(gatewayId.NatGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(gatewayId.NatGatewayName, natGatewayResourceName)
	locks.ByName(id.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(id.VirtualNetworkName, VirtualNetworkResourceName)

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
