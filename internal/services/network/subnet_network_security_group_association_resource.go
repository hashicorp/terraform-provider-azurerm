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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetNetworkSecurityGroupAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSubnetNetworkSecurityGroupAssociationCreate,
		Read:   resourceSubnetNetworkSecurityGroupAssociationRead,
		Delete: resourceSubnetNetworkSecurityGroupAssociationDelete,

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

			"network_security_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networksecuritygroups.ValidateNetworkSecurityGroupID,
			},
		},
	}
}

func resourceSubnetNetworkSecurityGroupAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.Subnets
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

	locks.ByName(networkSecurityGroupId.NetworkSecurityGroupName, networkSecurityGroupResourceName)
	defer locks.UnlockByName(networkSecurityGroupId.NetworkSecurityGroupName, networkSecurityGroupResourceName)

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
			if nsg := props.NetworkSecurityGroup; nsg != nil {
				// we're intentionally not checking the ID - if there's a NSG, it needs to be imported
				if nsg.Id != nil && model.Id != nil {
					return tf.ImportAsExistsAssociationError("azurerm_subnet_network_security_group_association", subnetId.ID(), *nsg.Id)
				}
			}

			props.NetworkSecurityGroup = &subnets.NetworkSecurityGroup{
				Id: pointer.To(networkSecurityGroupId.ID()),
			}
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *subnetId, *subnet.Model); err != nil {
		return fmt.Errorf("updating Network Security Group Association for %s: %+v", *subnetId, err)
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
		return fmt.Errorf("waiting for provisioning state of subnet for Network Security Group Association for %s: %+v", *subnetId, err)
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
		return fmt.Errorf("waiting for provisioning state of virtual network for Network Security Group Association for %s: %+v", *subnetId, err)
	}

	d.SetId(subnetId.ID())

	return resourceSubnetNetworkSecurityGroupAssociationRead(d, meta)
}

func resourceSubnetNetworkSecurityGroupAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] %s could not be found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("`model` was nil for %s", *id)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("`properties` was nil for %s", *id)
	}

	securityGroup := props.NetworkSecurityGroup
	if securityGroup == nil {
		log.Printf("[DEBUG] %s doesn't have a Network Security Group - removing from state!", *id)
		d.SetId("")
		return nil
	}

	d.Set("subnet_id", model.Id)
	d.Set("network_security_group_id", securityGroup.Id)

	return nil
}

func resourceSubnetNetworkSecurityGroupAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] %s could not be found - removing from state!", *id)
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := read.Model
	if model == nil {
		return fmt.Errorf("`model` was nil for %s", *id)
	}

	props := model.Properties
	if props == nil {
		return fmt.Errorf("`Properties` was nil for %s", *id)
	}

	if props.NetworkSecurityGroup == nil || props.NetworkSecurityGroup.Id == nil {
		log.Printf("[DEBUG] %s has no Network Security Group - removing from state!", *id)
		return nil
	}

	// once we have the network security group id to lock on, lock on that
	networkSecurityGroupId, err := networksecuritygroups.ParseNetworkSecurityGroupID(*props.NetworkSecurityGroup.Id)
	if err != nil {
		return err
	}

	locks.ByName(networkSecurityGroupId.NetworkSecurityGroupName, networkSecurityGroupResourceName)
	defer locks.UnlockByName(networkSecurityGroupId.NetworkSecurityGroupName, networkSecurityGroupResourceName)

	locks.ByName(id.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(id.VirtualNetworkName, VirtualNetworkResourceName)

	locks.ByName(id.SubnetName, SubnetResourceName)
	defer locks.UnlockByName(id.SubnetName, SubnetResourceName)

	// then re-retrieve it to ensure we've got the latest state
	read, err = client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] %s could not be found - removing from state!", *id)
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	read.Model.Properties.NetworkSecurityGroup = nil

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *read.Model); err != nil {
		return fmt.Errorf("removing Network Security Group Association from %s: %+v", *id, err)
	}

	return nil
}
