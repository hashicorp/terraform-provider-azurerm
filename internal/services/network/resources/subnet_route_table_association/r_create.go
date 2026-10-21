// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subnet_route_table_association

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetRouteTableAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Subnets
	vnetClient := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return err
	}

	routeTableId, err := routetables.ParseRouteTableID(d.Get("route_table_id").(string))
	if err != nil {
		return err
	}

	// find and lock chain of serialised change
	existingUnlocked, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existingUnlocked.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	var nsgIds []string
	var rtIds []string

	rtIds = append(rtIds, routeTableId.ID())

	if existingUnlocked.Model != nil && existingUnlocked.Model.Properties != nil {
		propsUnlocked := existingUnlocked.Model.Properties
		if propsUnlocked.NetworkSecurityGroup != nil && propsUnlocked.NetworkSecurityGroup.Id != nil {
			oldNsgId, err := networksecuritygroups.ParseNetworkSecurityGroupID(*propsUnlocked.NetworkSecurityGroup.Id)
			if err != nil {
				return fmt.Errorf("parsing existing Network Security Group ID: %+v", err)
			}
			nsgIds = append(nsgIds, oldNsgId.ID())
		}

		if propsUnlocked.RouteTable != nil && propsUnlocked.RouteTable.Id != nil {
			oldRtId, err := routetables.ParseRouteTableID(*propsUnlocked.RouteTable.Id)
			if err != nil {
				return fmt.Errorf("parsing existing Route Table ID: %+v", err)
			}
			rtIds = append(rtIds, oldRtId.ID())
		}
	}

	locks.MultipleByID(&nsgIds)
	defer locks.UnlockMultipleByID(&nsgIds)

	locks.MultipleByID(&rtIds)
	defer locks.UnlockMultipleByID(&rtIds)

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	locks.ByID(vnetId.ID())
	defer locks.UnlockByID(vnetId.ID())

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	// now we have exclusive access, we can read reliably for create
	subnet, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(subnet.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := subnet.Model; model != nil {
		if props := model.Properties; props != nil {
			if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if rt := props.RouteTable; rt != nil {
					// we're intentionally not checking the ID - if there's a RouteTable, it needs to be imported
					if rt.Id != nil && model.Id != nil {
						return tf.ImportAsExistsError("azurerm_subnet_route_table_association", *model.Id)
					}
				}
			}

			props.RouteTable = &subnets.RouteTable{
				Id: pointer.To(routeTableId.ID()),
			}
		}
	}

	// TODO: migrate this to a Composite ID
	if err := client.CreateOrUpdateCallbackThenPoll(ctx, *id, *subnet.Model, sdk.SetIDCallback(meta, id, d)); err != nil {
		return fmt.Errorf("updating Route Table Association for %s: %+v", id, err)
	}
	d.SetId(id.ID())

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    helpers.SubnetProvisioningStateRefreshFunc(ctx, client, *id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of Route Table Association for %s: %+v", id, err)
	}

	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    helpers.VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for Route Table Association for %s: %+v", id, err)
	}

	return resourceSubnetRouteTableAssociationRead(d, meta)
}
