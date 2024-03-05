// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
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
				ValidateFunc: validate.NetworkSecurityGroupID,
			},
		},
	}
}

func resourceSubnetNetworkSecurityGroupAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.Subnets
	vnetClient := meta.(*clients.Client).Network.VnetClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Subnet <-> Network Security Group Association creation.")

	subnetId := d.Get("subnet_id").(string)
	networkSecurityGroupId := d.Get("network_security_group_id").(string)

	parsedSubnetId, err := commonids.ParseSubnetID(subnetId)
	if err != nil {
		return err
	}

	parsedNetworkSecurityGroupId, err := parse.NetworkSecurityGroupID(networkSecurityGroupId)
	if err != nil {
		return err
	}

	locks.ByName(parsedNetworkSecurityGroupId.Name, networkSecurityGroupResourceName)
	defer locks.UnlockByName(parsedNetworkSecurityGroupId.Name, networkSecurityGroupResourceName)

	locks.ByName(parsedSubnetId.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(parsedSubnetId.VirtualNetworkName, VirtualNetworkResourceName)

	locks.ByName(parsedSubnetId.SubnetName, SubnetResourceName)
	defer locks.UnlockByName(parsedSubnetId.SubnetName, SubnetResourceName)

	subnet, err := client.Get(ctx, *parsedSubnetId, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(subnet.HttpResponse) {
			return fmt.Errorf("%s was not found", *parsedSubnetId)
		}

		return fmt.Errorf("retrieving %s: %+v", *parsedSubnetId, err)
	}

	if model := subnet.Model; model != nil {
		if props := model.Properties; props != nil {
			if nsg := props.NetworkSecurityGroup; nsg != nil {
				// we're intentionally not checking the ID - if there's a NSG, it needs to be imported
				if nsg.Id != nil && model.Id != nil {
					return tf.ImportAsExistsError("azurerm_subnet_network_security_group_association", *model.Id)
				}
			}

			props.NetworkSecurityGroup = &subnets.NetworkSecurityGroup{
				Id: utils.String(networkSecurityGroupId),
			}
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *parsedSubnetId, *subnet.Model); err != nil {
		return fmt.Errorf("updating Network Security Group Association for %s: %+v", *parsedSubnetId, err)
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
		return fmt.Errorf("waiting for provisioning state of subnet for Network Security Group Association for %s: %+v", *parsedSubnetId, err)
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
		return fmt.Errorf("waiting for provisioning state of virtual network for Network Security Group Association for %s: %+v", *parsedSubnetId, err)
	}

	d.SetId(parsedSubnetId.ID())

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
	parsedNetworkSecurityGroupId, err := parse.NetworkSecurityGroupID(*props.NetworkSecurityGroup.Id)
	if err != nil {
		return err
	}

	locks.ByName(parsedNetworkSecurityGroupId.Name, networkSecurityGroupResourceName)
	defer locks.UnlockByName(parsedNetworkSecurityGroupId.Name, networkSecurityGroupResourceName)

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
