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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
	vnetClient := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Subnet <-> Route Table Association creation.")

	id, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return err
	}

	routeTableId, err := routetables.ParseRouteTableID(d.Get("route_table_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(routeTableId.RouteTableName, routeTableResourceName)
	defer locks.UnlockByName(routeTableId.RouteTableName, routeTableResourceName)

	locks.ByName(id.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(id.VirtualNetworkName, VirtualNetworkResourceName)

	subnet, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(subnet.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
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
				Id: pointer.To(routeTableId.ID()),
			}
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *subnet.Model); err != nil {
		return fmt.Errorf("updating Route Table Association for %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    SubnetProvisioningStateRefreshFunc(ctx, client, *id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of Route Table Association for %s: %+v", id, err)
	}

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for Route Table Association for %s: %+v", id, err)
	}

	d.SetId(id.ID())

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
			log.Printf("[DEBUG] %s could not be found - removing from state!", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	routeTable := props.RouteTable
	if routeTable == nil {
		log.Printf("[DEBUG] %s doesn't have a Route Table - removing from state!", id)
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
			log.Printf("[DEBUG] %s could not be found - removing from state!", id)
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := read.Model
	if model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if props.RouteTable == nil || props.RouteTable.Id == nil {
		log.Printf("[DEBUG] %s has no Route Table - removing from state!", id)
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
			log.Printf("[DEBUG] %s could not be found - removing from state!", id)
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	read.Model.Properties.RouteTable = nil

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *read.Model); err != nil {
		return fmt.Errorf("removing Route Table Association from %s: %+v", id, err)
	}

	return nil
}
