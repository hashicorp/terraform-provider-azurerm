// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subnet_network_security_group_association

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

func resourceSubnetNetworkSecurityGroupAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Subnets
	vnetClient := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subnetId, err := commonids.ParseSubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return err
	}

	networkSecurityGroupId, err := networksecuritygroups.ParseNetworkSecurityGroupID(d.Get("network_security_group_id").(string))
	if err != nil {
		return err
	}

	// find and lock chain of sequential change
	existingUnlocked, err := client.Get(ctx, *subnetId, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existingUnlocked.HttpResponse) {
			return fmt.Errorf("%s was not found", *subnetId)
		}
		return fmt.Errorf("retrieving %s: %+v", *subnetId, err)
	}

	var nsgIds []string
	var rtIds []string

	nsgIds = append(nsgIds, networkSecurityGroupId.ID())

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
			rtId, err := routetables.ParseRouteTableID(*propsUnlocked.RouteTable.Id)
			if err != nil {
				return fmt.Errorf("parsing existing Route Table ID: %+v", err)
			}
			rtIds = append(rtIds, rtId.ID())
		}
	}

	locks.MultipleByID(&nsgIds)
	defer locks.UnlockMultipleByID(&nsgIds)

	locks.MultipleByID(&rtIds)
	defer locks.UnlockMultipleByID(&rtIds)

	vnetId := commonids.NewVirtualNetworkID(subnetId.SubscriptionId, subnetId.ResourceGroupName, subnetId.VirtualNetworkName)
	locks.ByID(vnetId.ID())
	defer locks.UnlockByID(vnetId.ID())

	locks.ByID(subnetId.ID())
	defer locks.UnlockByID(subnetId.ID())

	// Now we have exclusive access, we can read reliably
	subnet, err := client.Get(ctx, *subnetId, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(subnet.HttpResponse) {
			return fmt.Errorf("%s was not found", *subnetId)
		}
		return fmt.Errorf("retrieving %s: %+v", *subnetId, err)
	}

	if model := subnet.Model; model != nil {
		if props := model.Properties; props != nil {
			if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if nsg := props.NetworkSecurityGroup; nsg != nil {
					// we're intentionally not checking the ID - if there's a NSG, it needs to be imported
					if nsg.Id != nil && model.Id != nil {
						return tf.ImportAsExistsAssociationError("azurerm_subnet_network_security_group_association", subnetId.ID(), *nsg.Id)
					}
				}
			}

			props.NetworkSecurityGroup = &subnets.NetworkSecurityGroup{
				Id: pointer.To(networkSecurityGroupId.ID()),
			}
		}
	}

	// TODO: migrate to a composite ID
	if err := client.CreateOrUpdateCallbackThenPoll(ctx, *subnetId, *subnet.Model, sdk.SetIDCallback(meta, subnetId, d)); err != nil {
		return fmt.Errorf("updating Network Security Group Association for %s: %+v", *subnetId, err)
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
		return fmt.Errorf("waiting for provisioning state of subnet for Network Security Group Association for %s: %+v", *subnetId, err)
	}

	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    helpers.VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for Network Security Group Association for %s: %+v", *subnetId, err)
	}

	return resourceSubnetNetworkSecurityGroupAssociationRead(d, meta)
}
