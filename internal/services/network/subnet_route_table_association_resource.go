// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceSubnetRouteTableAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSubnetRouteTableAssociationCreate,
		Read:   resourceSubnetRouteTableAssociationRead,
		Delete: resourceSubnetRouteTableAssociationDelete,

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

			"route_table_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: routetables.ValidateRouteTableID,
			},
		},
	}
}

func resourceSubnetRouteTableAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.Subnets
	vnetClient := meta.(*clients.Client).Network.VnetClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Subnet <-> Route Table Association creation.")

	subnetId := d.Get("subnet_id").(string)
	routeTableId := d.Get("route_table_id").(string)

	parsedSubnetId, err := commonids.ParseSubnetID(subnetId)
	if err != nil {
		return err
	}

	parsedRouteTableId, err := routetables.ParseRouteTableID(routeTableId)
	if err != nil {
		return err
	}

	locks.ByName(parsedRouteTableId.RouteTableName, routeTableResourceName)
	defer locks.UnlockByName(parsedRouteTableId.RouteTableName, routeTableResourceName)

	locks.ByName(parsedSubnetId.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(parsedSubnetId.VirtualNetworkName, VirtualNetworkResourceName)

	subnet, err := client.Get(ctx, *parsedSubnetId, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(subnet.HttpResponse) {
			return fmt.Errorf("Subnet %q (Virtual Network %q / Resource Group %q) was not found!", parsedSubnetId.SubnetName, parsedSubnetId.VirtualNetworkName, parsedSubnetId.ResourceGroupName)
		}

		return fmt.Errorf("retrieving Subnet %q (Virtual Network %q / Resource Group %q): %+v", parsedSubnetId.SubnetName, parsedSubnetId.VirtualNetworkName, parsedSubnetId.ResourceGroupName, err)
	}

	if model := subnet.Model; model != nil {
		if props := model.Properties; props != nil {
			if rt := props.RouteTable; rt != nil {
				// we're intentionally not checking the ID - if there's a RouteTable, it needs to be imported
				if rt.Id != nil && model.Id != nil {
					return tf.ImportAsExistsError("azurerm_subnet_route_table_association", *model.Id)
				}
			}

			props.RouteTable = &subnets.RouteTable{
				Id: utils.String(routeTableId),
			}
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *parsedSubnetId, *subnet.Model); err != nil {
		return fmt.Errorf("updating Route Table Association for Subnet %q (Virtual Network %q / Resource Group %q): %+v", parsedSubnetId.SubnetName, parsedSubnetId.VirtualNetworkName, parsedSubnetId.ResourceGroupName, err)
	}

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(network.ProvisioningStateUpdating)},
		Target:     []string{string(network.ProvisioningStateSucceeded)},
		Refresh:    SubnetProvisioningStateRefreshFunc(ctx, client, *parsedSubnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of subnet for Route Table Association for Subnet %q (Virtual Network %q / Resource Group %q): %+v", parsedSubnetId.SubnetName, parsedSubnetId.VirtualNetworkName, parsedSubnetId.ResourceGroupName, err)
	}

	vnetId := commonids.NewVirtualNetworkID(parsedSubnetId.SubscriptionId, parsedSubnetId.ResourceGroupName, parsedSubnetId.VirtualNetworkName)
	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(network.ProvisioningStateUpdating)},
		Target:     []string{string(network.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for Route Table Association for Subnet %q (Virtual Network %q / Resource Group %q): %+v", parsedSubnetId.SubnetName, parsedSubnetId.VirtualNetworkName, parsedSubnetId.ResourceGroupName, err)
	}

	d.SetId(parsedSubnetId.ID())

	return resourceSubnetRouteTableAssociationRead(d, meta)
}

func resourceSubnetRouteTableAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.Subnets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Subnet %q (Virtual Network %q / Resource Group %q) could not be found - removing from state!", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Subnet %q (Virtual Network %q / Resource Group %q): %+v", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("Error: `model` was nil for Subnet %q (Virtual Network %q / Resource Group %q)", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("Error: `properties` was nil for Subnet %q (Virtual Network %q / Resource Group %q)", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName)
	}

	routeTable := props.RouteTable
	if routeTable == nil {
		log.Printf("[DEBUG] Subnet %q (Virtual Network %q / Resource Group %q) doesn't have a Route Table - removing from state!", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName)
		d.SetId("")
		return nil
	}

	d.Set("subnet_id", model.Id)
	d.Set("route_table_id", routeTable.Id)

	return nil
}

func resourceSubnetRouteTableAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.Subnets
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	// retrieve the subnet
	read, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] Subnet %q (Virtual Network %q / Resource Group %q) could not be found - removing from state!", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName)
			return nil
		}

		return fmt.Errorf("retrieving Subnet %q (Virtual Network %q / Resource Group %q): %+v", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName, err)
	}

	model := read.Model
	if model == nil {
		return fmt.Errorf("`model` was nil for Subnet %q (Virtual Network %q / Resource Group %q)", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("`Properties` was nil for Subnet %q (Virtual Network %q / Resource Group %q)", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName)
	}

	if props.RouteTable == nil || props.RouteTable.Id == nil {
		log.Printf("[DEBUG] Subnet %q (Virtual Network %q / Resource Group %q) has no Route Table - removing from state!", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName)
		return nil
	}

	// once we have the route table id to lock on, lock on that
	parsedRouteTableId, err := routetables.ParseRouteTableID(*props.RouteTable.Id)
	if err != nil {
		return err
	}

	locks.ByName(parsedRouteTableId.RouteTableName, routeTableResourceName)
	defer locks.UnlockByName(parsedRouteTableId.RouteTableName, routeTableResourceName)

	locks.ByName(id.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(id.VirtualNetworkName, VirtualNetworkResourceName)

	// then re-retrieve it to ensure we've got the latest state
	read, err = client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] Subnet %q (Virtual Network %q / Resource Group %q) could not be found - removing from state!", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName)
			return nil
		}

		return fmt.Errorf("retrieving Subnet %q (Virtual Network %q / Resource Group %q): %+v", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName, err)
	}

	read.Model.Properties.RouteTable = nil

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *read.Model); err != nil {
		return fmt.Errorf("removing Route Table Association from Subnet %q (Virtual Network %q / Resource Group %q): %+v", id.SubnetName, id.VirtualNetworkName, id.ResourceGroupName, err)
	}

	return nil
}
