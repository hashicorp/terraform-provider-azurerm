// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetworkInterfaceApplicationSecurityGroupAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceApplicationSecurityGroupAssociationCreate,
		Read:   resourceNetworkInterfaceApplicationSecurityGroupAssociationRead,
		Delete: resourceNetworkInterfaceApplicationSecurityGroupAssociationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseCompositeResourceID(id, &commonids.NetworkInterfaceId{}, &applicationsecuritygroups.ApplicationSecurityGroupId{})
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NetworkInterfaceApplicationSecurityGroupAssociationV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"network_interface_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateNetworkInterfaceID,
			},

			"application_security_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: applicationsecuritygroups.ValidateApplicationSecurityGroupID,
			},
		},
	}
}

func resourceNetworkInterfaceApplicationSecurityGroupAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Network Interface <-> Application Security Group Association creation.")

	applicationSecurityGroupId, err := applicationsecuritygroups.ParseApplicationSecurityGroupID(d.Get("application_security_group_id").(string))
	if err != nil {
		return err
	}
	networkInterfaceId, err := commonids.ParseNetworkInterfaceID(d.Get("network_interface_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(networkInterfaceId.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(networkInterfaceId.NetworkInterfaceName, networkInterfaceResourceName)

	read, err := client.Get(ctx, *networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", networkInterfaceId.ID())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *networkInterfaceId, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", networkInterfaceId)
	}

	props := read.Model.Properties
	if props == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", networkInterfaceId)
	}
	if props.IPConfigurations == nil {
		return fmt.Errorf("retrieving %s: `properties.ipConfigurations` was nil", networkInterfaceId)
	}

	info := parseFieldsFromNetworkInterface(*props)
	id := commonids.NewCompositeResourceID(networkInterfaceId, applicationSecurityGroupId)

	if utils.SliceContainsValue(info.applicationSecurityGroupIDs, applicationSecurityGroupId.ID()) {
		return tf.ImportAsExistsError("azurerm_network_interface_application_security_group_association", id.ID())
	}

	info.applicationSecurityGroupIDs = append(info.applicationSecurityGroupIDs, applicationSecurityGroupId.ID())

	props.IPConfigurations = mapFieldsToNetworkInterface(props.IPConfigurations, info)

	err = client.CreateOrUpdateThenPoll(ctx, *networkInterfaceId, *read.Model)
	if err != nil {
		return fmt.Errorf("updating Application Security Group Association for %s: %+v", *networkInterfaceId, err)
	}

	d.SetId(id.ID())

	return resourceNetworkInterfaceApplicationSecurityGroupAssociationRead(d, meta)
}

func resourceNetworkInterfaceApplicationSecurityGroupAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceId{}, &applicationsecuritygroups.ApplicationSecurityGroupId{})
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id.First, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id.First)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if model := read.Model; model != nil {
		nicProps := read.Model.Properties
		if nicProps == nil {
			return fmt.Errorf("retrieving %s: `properties` was nil", id.First)
		}

		info := parseFieldsFromNetworkInterface(*nicProps)
		exists := false
		for _, groupId := range info.applicationSecurityGroupIDs {
			if groupId == id.Second.ID() {
				exists = true
			}
		}

		if !exists {
			log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.First, id.Second)
			d.SetId("")
			return nil
		}
	}

	d.Set("application_security_group_id", id.Second.ID())
	d.Set("network_interface_id", id.First.ID())

	return nil
}

func resourceNetworkInterfaceApplicationSecurityGroupAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceId{}, &applicationsecuritygroups.ApplicationSecurityGroupId{})
	if err != nil {
		return err
	}

	locks.ByName(id.First.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(id.First.NetworkInterfaceName, networkInterfaceResourceName)

	read, err := client.Get(ctx, *id.First, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			return fmt.Errorf("%s was not found", id.First)
		}

		return fmt.Errorf("retrieving %s: %+v", id.First, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id.First)
	}

	props := read.Model.Properties
	if props == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id.First)
	}

	if props.IPConfigurations == nil {
		return fmt.Errorf("retrieving %s: `properties.ipConfigurations` was nil)", id.First)
	}

	info := parseFieldsFromNetworkInterface(*props)

	applicationSecurityGroupIds := make([]string, 0)
	for _, v := range info.applicationSecurityGroupIDs {
		if v != id.Second.ID() {
			applicationSecurityGroupIds = append(applicationSecurityGroupIds, v)
		}
	}
	info.applicationSecurityGroupIDs = applicationSecurityGroupIds
	props.IPConfigurations = mapFieldsToNetworkInterface(props.IPConfigurations, info)

	err = client.CreateOrUpdateThenPoll(ctx, *id.First, *read.Model)
	if err != nil {
		return fmt.Errorf("removing Application Security Group for %s: %+v", id.First, err)
	}

	return nil
}
