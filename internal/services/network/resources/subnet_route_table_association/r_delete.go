// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet_route_table_association

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetRouteTableAssociationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Subnets
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	// find and lock chain of serialised change
	readUnlocked, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(readUnlocked.HttpResponse) {
			log.Printf("[DEBUG] %s could not be found - removing from state!", id)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if readUnlocked.Model == nil || readUnlocked.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model` or `properties` was nil", id)
	}

	propsUnlocked := readUnlocked.Model.Properties

	if propsUnlocked.RouteTable == nil || propsUnlocked.RouteTable.Id == nil {
		log.Printf("[DEBUG] %s has no Route Table - removing from state!", id)
		return nil
	}

	var nsgIds []string
	var rtIds []string

	if propsUnlocked.NetworkSecurityGroup != nil && propsUnlocked.NetworkSecurityGroup.Id != nil {
		nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupID(*propsUnlocked.NetworkSecurityGroup.Id)
		if err != nil {
			return err
		}
		nsgIds = append(nsgIds, nsgId.ID())
	}

	parsedRouteTableId, err := routetables.ParseRouteTableID(*propsUnlocked.RouteTable.Id)
	if err != nil {
		return err
	}
	rtIds = append(rtIds, parsedRouteTableId.ID())

	locks.MultipleByID(&nsgIds)
	defer locks.UnlockMultipleByID(&nsgIds)

	locks.MultipleByID(&rtIds)
	defer locks.UnlockMultipleByID(&rtIds)

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	locks.ByID(vnetId.ID())
	defer locks.UnlockByID(vnetId.ID())

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	// Now we have the locks, we can try to delete
	read, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
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
